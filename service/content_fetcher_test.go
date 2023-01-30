package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedContentRepository struct {
	mock.Mock
}

func (m *MockedContentRepository) FindByID(ctx context.Context, url string) (*domain.ObserverTrace, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(*domain.ObserverTrace), args.Error(1)
}
func (m *MockedContentRepository) Save(ctx context.Context, content domain.ObserverTrace) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

func TestFetchAndCompareWithIdenticalStringsThenResultEqualsUnchanged(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("<html><body><span>Out of Stock</span>My content</body></html>"))
	}))

	data, _ := base64.StdEncoding.Strict().DecodeString("H4sIAAAAAAAA/7LJKMnNsbNJyk+ptLMpLkjMs/MvLVHIT1MILslPzrbRBwv5Viok5+eVpOaV2OhDVOqDtQECAAD//1Jq3a09AAAA")
	dbContent := &domain.ObserverTrace{
		Data: data,
		URL:  fmt.Sprintf("%s/my-endpoint", mockServer.URL),
	}
	mockedTrackingRepository := new(MockedContentRepository)
	mockedTrackingRepository.On("FindByID", mock.Anything, mock.Anything).Return(dbContent, nil)
	mockedTrackingRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	httpClient := &http.Client{}
	contentFetcher := NewContentFetcher(mockedTrackingRepository, httpClient)

	// Act
	result, err := contentFetcher.FetchAndCompare(context.TODO(), fmt.Sprintf("%s/my-endpoint", mockServer.URL), "body > span")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, domain.Unchanged, result.State)

}

func TestFetchAndCompareWithDifferentStringsThenResultEqualsUpdated(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("<html><body><span>In Stock</span>My content</body></html>"))
	}))

	data, _ := base64.StdEncoding.Strict().DecodeString("H4sIAAAAAAAA/7LJKMnNsbNJyk+ptLMpLkjMs/MvLVHIT1MILslPzrbRBwv5Viok5+eVpOaV2OhDVOqDtQECAAD//1Jq3a09AAAA")
	dbContent := &domain.ObserverTrace{
		Data: data,
		URL:  fmt.Sprintf("%s/my-endpoint", mockServer.URL),
	}
	mockedTrackingRepository := new(MockedContentRepository)
	mockedTrackingRepository.On("FindByID", mock.Anything, mock.Anything).Return(dbContent, nil)
	mockedTrackingRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	httpClient := &http.Client{}
	contentFetcher := NewContentFetcher(mockedTrackingRepository, httpClient)

	// Act
	result, err := contentFetcher.FetchAndCompare(context.TODO(), fmt.Sprintf("%s/my-endpoint", mockServer.URL), "body > span")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, domain.Updated, result.State)

}
