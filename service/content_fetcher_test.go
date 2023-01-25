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

func (m *MockedContentRepository) FindByID(ctx context.Context, url string) (*domain.Content, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(*domain.Content), args.Error(1)
}

func TestFetchAndCompareWithIdenticalStringsThenResultEqualsUnchanged(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("<html><body>My content</body></html>"))
	}))

	dbContent := &domain.Content{
		Data:     "<html><body>My content</body></html>",
		URL:      fmt.Sprintf("%s/my-endpoint", mockServer.URL),
		IsActive: true,
	}
	mockedTrackingRepository := new(MockedContentRepository)
	mockedTrackingRepository.On("FindByID", mock.Anything, mock.Anything).Return(dbContent, nil)

	httpClient := &http.Client{}
	contentFetcher := NewContentFetcher(mockedTrackingRepository, httpClient)

	// Act
	result, err := contentFetcher.FetchAndCompare(context.TODO(), fmt.Sprintf("%s/my-endpoint", mockServer.URL))

	// Assert
	expectedValue := domain.Unchanged
	assert.Nil(t, err)
	assert.Equal(t, &expectedValue, result)

}

func TestFetchAndCompareWithDifferentStringsThenResultEqualsUpdated(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("<html><body><span>In Stock</span>My content</body></html>"))
	}))

	dbContent := &domain.Content{
		Data:     "<html><body><span>Out of stock</span>My content</body></html>",
		URL:      fmt.Sprintf("%s/my-endpoint", mockServer.URL),
		IsActive: true,
	}
	mockedTrackingRepository := new(MockedContentRepository)
	mockedTrackingRepository.On("FindByID", mock.Anything, mock.Anything).Return(dbContent, nil)

	httpClient := &http.Client{}
	contentFetcher := NewContentFetcher(mockedTrackingRepository, httpClient)

	// Act
	result, err := contentFetcher.FetchAndCompare(context.TODO(), fmt.Sprintf("%s/my-endpoint", mockServer.URL))

	// Assert
	expectedValue := domain.Updated
	assert.Nil(t, err)
	assert.Equal(t, &expectedValue, result)

}
