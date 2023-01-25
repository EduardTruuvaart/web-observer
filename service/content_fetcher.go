package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service/compressor"
	"golang.org/x/net/html"
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

	previousHTML, _ := html.Parse(strings.NewReader(string(decompressedData)))
	latestHTML, _ := html.Parse(strings.NewReader(data))

	if reflect.DeepEqual(previousHTML, latestHTML) {
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
