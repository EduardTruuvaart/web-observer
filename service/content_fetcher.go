package service

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/EduardTruuvaart/web-observer/domain"
	htmlcompareresult "github.com/EduardTruuvaart/web-observer/domain/htmlcompare"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service/compressor"
	"github.com/EduardTruuvaart/web-observer/service/htmldiff"
)

type ContentFetcher struct {
	contentRepository repository.ContentRepository
	httpClient        *http.Client
}

func NewContentFetcher(contentRepository repository.ContentRepository, httpClient *http.Client) *ContentFetcher {
	return &ContentFetcher{
		contentRepository: contentRepository,
		httpClient:        httpClient,
	}
}

func (c *ContentFetcher) FetchAndCompare(ctx context.Context, url string, cssSelector string) (domain.FetchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return domain.FetchResult{}, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		fmt.Printf("Got error calling httpClient.Do: %s\n", err)
		return domain.FetchResult{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := string(body)

	if err != nil {
		fmt.Printf("Got error calling ioutil.ReadAll: %s\n", err)
		return domain.FetchResult{}, err
	}

	savedResult, err := c.contentRepository.FindByID(ctx, url)
	if err != nil {
		fmt.Printf("Got error calling contentRepository.FindByID: %s\n", err)
		return domain.FetchResult{}, err
	}

	if savedResult == nil {
		err = c.saveLatestContent(ctx, url, data, cssSelector, true)

		if err != nil {
			fmt.Printf("Got error calling saveLatestContent: %s\n", err)
			return domain.FetchResult{}, err
		}

		return domain.FetchResult{
			State: domain.NewContentIsAdded,
		}, nil
	}

	decompressedData, err := compressor.Decompress(savedResult.Data)

	if err != nil {
		fmt.Printf("Got error calling Decompress: %s\n", err)
		return domain.FetchResult{}, err
	}

	escapedData := html.EscapeString(data)
	escapedPreviousData := html.EscapeString(string(decompressedData))

	if strings.Compare(escapedData, escapedPreviousData) == 0 {
		return domain.FetchResult{
			State: domain.Unchanged,
		}, nil
	}

	err = c.saveLatestContent(ctx, url, data, cssSelector, true)

	if err != nil {
		fmt.Printf("Got error calling saveLatestContent 2: %s\n", err)
		return domain.FetchResult{
			State: domain.Updated,
		}, err
	}

	result, err := htmldiff.CompareDocumentSection(string(decompressedData), data, "body")

	if err != nil {
		fmt.Printf("Got error comparing html documents: %s\n", err)
		return domain.FetchResult{
			State: domain.Updated,
		}, err
	}

	diffString := ""
	if result.State == htmlcompareresult.SelectionNotFoundInSource {
		diffString = "Selection not found in source"
	}

	if result.State == htmlcompareresult.SelectionNotFoundInTarget {
		diffString = "Selection not found in target"
	}

	if result.State == htmlcompareresult.Different {
		for _, diff := range result.Differences {
			diffString += fmt.Sprintf("%s\n", diff)
		}
	}

	return domain.FetchResult{
		State:      domain.Updated,
		Difference: diffString,
		DiffSize:   result.DiffSize,
	}, nil
}

func (c *ContentFetcher) saveLatestContent(ctx context.Context, url, data, cssSelector string, isActive bool) error {
	compressedData, err := compressor.Compress([]byte(data))

	//fmt.Printf("Compressed data size: %dkb\n", len(compressedData)/1024)

	if err != nil {
		return err
	}

	content := domain.Content{
		URL:         url,
		Data:        compressedData,
		CssSelector: cssSelector,
		IsActive:    isActive,
	}
	err = c.contentRepository.Save(ctx, content)

	return err
}
