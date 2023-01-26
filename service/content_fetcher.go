package service

import (
	"context"
	"html"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service/compressor"
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
		return domain.Unchanged, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return domain.Unchanged, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := string(body)

	if err != nil {
		return domain.Unchanged, err
	}

	savedResult, err := c.contentRepository.FindByID(ctx, url)
	if err != nil {
		return domain.Unchanged, err
	}

	if savedResult == nil {
		err = c.saveLatestContent(ctx, url, data, true)

		if err != nil {
			return domain.Unchanged, err
		}

		result := domain.NewContentIsAdded
		return result, nil
	}

	decompressedData, err := compressor.Decompress(savedResult.Data)

	if err != nil {
		return domain.Unchanged, err
	}

	escapedData := html.EscapeString(data)
	escapedPreviousData := html.EscapeString(string(decompressedData))

	if strings.Compare(escapedData, escapedPreviousData) == 0 {
		result := domain.Unchanged
		return result, nil
	}

	err = c.saveLatestContent(ctx, url, data, true)

	if err != nil {
		return domain.Updated, err
	}

	result := domain.Updated
	return result, nil
}

func (c *ContentFetcher) saveLatestContent(ctx context.Context, url string, data string, isActive bool) error {
	compressedData, err := compressor.Compress([]byte(data))

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
