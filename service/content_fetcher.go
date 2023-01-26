package service

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service/compressor"
	"github.com/sergi/go-diff/diffmatchpatch"
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

func (c *ContentFetcher) FetchAndCompare(ctx context.Context, url string) (domain.FetchResult, error) {
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
		err = c.saveLatestContent(ctx, url, data, true)

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

	err = c.saveLatestContent(ctx, url, data, true)

	if err != nil {
		fmt.Printf("Got error calling saveLatestContent 2: %s\n", err)
		return domain.FetchResult{
			State: domain.Updated,
		}, err
	}

	// var cfg = &htmldiff.Config{
	// 	Granularity:  5,
	// 	InsertedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: palegreen;"}},
	// 	DeletedSpan:  []htmldiff.Attribute{{Key: "style", Val: "background-color: lightpink;"}},
	// 	ReplacedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: lightskyblue;"}},
	// 	CleanTags:    []string{""},
	// }
	// res, err := cfg.HTMLdiff([]string{string(decompressedData), data})

	// if err != nil {
	// 	fmt.Printf("Got error calling HTMLdiff: %s\n", err)
	// 	return domain.FetchResult{
	// 		State: domain.Updated,
	// 	}, err
	// }

	// mergedHTML := res[0]
	// fmt.Println(mergedHTML)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(decompressedData), data, false)
	fmt.Println(len(diffs))
	diffString := ""
	for _, diff := range diffs {
		diffString += fmt.Sprintf("%s\n", diff.Text)
	}

	return domain.FetchResult{
		State:      domain.Updated,
		Difference: diffString,
		DiffSize:   len(diffs),
	}, nil
}

func (c *ContentFetcher) saveLatestContent(ctx context.Context, url string, data string, isActive bool) error {
	compressedData, err := compressor.Compress([]byte(data))

	//fmt.Printf("Compressed data size: %dkb\n", len(compressedData)/1024)

	if err != nil {
		return err
	}

	content := domain.Content{
		URL:      url,
		Data:     compressedData,
		IsActive: isActive,
	}
	err = c.contentRepository.Save(ctx, content)

	return err
}
