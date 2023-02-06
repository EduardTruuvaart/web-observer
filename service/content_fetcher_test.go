package service

import (
	"context"
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

func (m *MockedContentRepository) Create(ctx context.Context, chatID int64) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}
func (m *MockedContentRepository) FindByID(ctx context.Context, chatID int64) (*domain.ObserverTrace, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).(*domain.ObserverTrace), args.Error(1)
}
func (m *MockedContentRepository) UpdateWithData(ctx context.Context, chatID int64, url string, data []byte) error {
	args := m.Called(ctx, chatID, url, data)
	return args.Error(0)
}
func (m *MockedContentRepository) UpdateWithUrl(ctx context.Context, chatID int64, url string) error {
	args := m.Called(ctx, chatID, url)
	return args.Error(0)
}
func (m *MockedContentRepository) UpdateWithSelectorAndActivate(ctx context.Context, chatID int64, cssSelector string) error {
	args := m.Called(ctx, chatID, cssSelector)
	return args.Error(0)
}
func (m *MockedContentRepository) Delete(ctx context.Context, chatID int64) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func TestFetchAndCompareWithIdenticalStringsThenResultEqualsUnchanged(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("<html><body><span>Out of Stock</span>My content</body></html>"))
	}))

	data := []byte("<html><body><span>Out of Stock</span>My content</body></html>")
	url := fmt.Sprintf("%s/my-endpoint", mockServer.URL)
	dbContent := &domain.ObserverTrace{
		Data: &data,
		URL:  &url,
	}
	mockedTrackingRepository := new(MockedContentRepository)
	mockedTrackingRepository.On("FindByID", mock.Anything, mock.Anything).Return(dbContent, nil)
	mockedTrackingRepository.On("UpdateWithData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	httpClient := &http.Client{}
	contentFetcher := NewContentFetcher(mockedTrackingRepository, httpClient)

	// Act
	result, err := contentFetcher.FetchAndCompare(context.TODO(), 123, fmt.Sprintf("%s/my-endpoint", mockServer.URL), "body > span")

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

	data := []byte("<html><body><span>Out of Stock</span>My content</body></html>")

	url := fmt.Sprintf("%s/my-endpoint", mockServer.URL)
	dbContent := &domain.ObserverTrace{
		Data: &data,
		URL:  &url,
	}
	mockedTrackingRepository := new(MockedContentRepository)
	mockedTrackingRepository.On("FindByID", mock.Anything, mock.Anything).Return(dbContent, nil)
	mockedTrackingRepository.On("UpdateWithData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	httpClient := &http.Client{}
	contentFetcher := NewContentFetcher(mockedTrackingRepository, httpClient)

	// Act
	result, err := contentFetcher.FetchAndCompare(context.TODO(), 123, fmt.Sprintf("%s/my-endpoint", mockServer.URL), "body > span")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, domain.Updated, result.State)

}
