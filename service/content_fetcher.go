package service

import (
	"net/url"

	"github.com/EduardTruuvaart/web-observer/repository"
)

type ContentFetcher struct {
	contentRepository repository.ContentRepository
}

func NewContentFetcher(contentRepository repository.ContentRepository) *ContentFetcher {
	return &ContentFetcher{
		contentRepository: contentRepository,
	}
}

func (f *ContentFetcher) FetchAndCompare(url url.URL) error {
	panic("not implemented")
}
