package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/EduardTruuvaart/web-observer/domain"
	htmlcompareresult "github.com/EduardTruuvaart/web-observer/domain/htmlcompare"
	"github.com/EduardTruuvaart/web-observer/repository"
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

func (c *ContentFetcher) FetchAndCompare(ctx context.Context, chatID int64, url string, cssSelector string) (domain.FetchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-GB,en;q=0.9,ru;q=0.8")

	if err != nil {
		return domain.FetchResult{}, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		fmt.Printf("Got error calling httpClient.Do: %s\n", err)
		return domain.FetchResult{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	data := string(body)

	//log.Printf("Response: %s\n", body)

	if err != nil {
		fmt.Printf("Got error calling ioutil.ReadAll: %s\n", err)
		return domain.FetchResult{}, err
	}

	savedResult, err := c.contentRepository.FindByID(ctx, chatID, url)
	if err != nil {
		fmt.Printf("Got error calling contentRepository.FindByID: %s\n", err)
		return domain.FetchResult{}, err
	}

	err = c.contentRepository.UpdateWithData(ctx, chatID, url, []byte(data))

	if err != nil {
		fmt.Printf("Got error calling saveLatestContent: %s\n", err)
		return domain.FetchResult{}, err
	}

	if savedResult == nil {
		return domain.FetchResult{
			State: domain.NewContentIsAdded,
		}, nil
	}

	result, err := htmldiff.CompareDocumentSection(string(*savedResult.Data), data, cssSelector)

	if err != nil {
		fmt.Printf("Got error comparing html documents: %s\n", err)
		return domain.FetchResult{
			State: domain.Updated,
		}, err
	}

	if result.State == htmlcompareresult.Identical {
		return domain.FetchResult{
			State: domain.Unchanged,
		}, nil
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
