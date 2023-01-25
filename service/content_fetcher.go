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

func (c *ContentFetcher) FetchAndCompare(ctx context.Context, url string) (*domain.FetchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := string(body)

	if err != nil {
		return nil, err
	}

	savedResult, err := c.contentRepository.FindByID(ctx, url)
	if err != nil {
		return nil, err
	}

	if savedResult == nil {
		compressedData, err := compressor.Compress([]byte(data))

		if err != nil {
			return nil, err
		}

		content := domain.Content{
			URL:      url,
			Data:     compressedData,
			IsActive: true,
		}
		err = c.contentRepository.Save(ctx, content)

		if err != nil {
			return nil, err
		}

		result := domain.NewContentIsAdded
		return &result, nil
	}

	decompressedData, err := compressor.Decompress(savedResult.Data)

	if err != nil {
		return nil, err
	}

	previousHTML, _ := html.Parse(strings.NewReader(string(decompressedData)))
	latestHTML, _ := html.Parse(strings.NewReader(data))

	if reflect.DeepEqual(previousHTML, latestHTML) {
		result := domain.Unchanged
		return &result, nil
	}

	result := domain.Updated
	return &result, nil
}
