package service

import (
	"context"
	"errors"
	"github.com/fir1/news/internal/news/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockService is a mock implementation of the Service struct that satisfies the NewsFetcher interface
type MockService struct {
	mock.Mock
}

// fetchNews is the mocked implementation of the fetchNews function
func (m *MockService) fetchNews(ctx context.Context, feedURL string) (RSS, error) {
	// Declare the arguments that the method should receive
	args := m.Called(ctx, feedURL)

	// Extract and return the results from the arguments
	return args.Get(0).(RSS), args.Error(1)
}

func TestListNews(t *testing.T) {
	// Sample test data
	ctx := context.Background()
	providers := []model.NewsProvider{model.NewsProviderBBC}
	categories := []string{"technology"}

	// Create the mock service
	mockService := new(MockService)

	// Set the expected method calls and return values for fetchNews
	rssData := RSS{
		Channel: Channel{
			Image: Image{
				URL: "https://www.example.com/image.jpg",
			},
			Items: []Item{
				{
					Title:       CDATA{Text: "Item 1 Title"},
					Link:        "https://www.example.com/item1",
					Description: "This is the description of item 1.",
					PubDate:     "Mon, 04 Jan 2023 15:04:05 GMT",
				},
				{
					Title:       CDATA{Text: "Item 2 Title"},
					Link:        "https://www.example.com/item2",
					Description: "This is the description of item 2.",
					PubDate:     "Mon, 02 Jan 2023 15:04:05 GMT",
				},
			},
		},
	}

	mockService.On("fetchNews", ctx, "http://feeds.bbci.co.uk/news/technology/rss.xml").Return(rssData, nil)

	// Create the Service instance using the mockService
	service := Service{
		NewsFetcher: mockService,
	}

	// Create test cases using table-driven testing
	testCases := []struct {
		name           string
		listNewsParams ListNewsParams
		expectedError  error
	}{
		{
			name: "Success",
			listNewsParams: ListNewsParams{
				Categories:        &categories,
				Providers:         &providers,
				NewsSourceURL:     nil,
				SortByPublishDate: "",
			},
			expectedError: nil,
		},
		{
			name: "Failure",
			listNewsParams: ListNewsParams{
				Categories:        nil,
				Providers:         &providers,
				NewsSourceURL:     StrPointer("https://www.bbc.co.uk/news/uk-politics-66268541?at_medium=RSS&amp;at_campaign=KARANGA"),
				SortByPublishDate: "",
			},
			expectedError: errors.New("please provide one of value for providers or news_source_url can not proceed both"),
		},
		{
			name: "Failure",
			listNewsParams: ListNewsParams{
				Categories:        nil,
				Providers:         &providers,
				SortByPublishDate: "random sort",
			},
			expectedError: errors.New("please provide a valid sort by publish date ASC or DESC"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function under test
			getNewsResponse, err := service.ListNews(ctx, tc.listNewsParams)
			if tc.name == "Success" {
				// Assertions
				assert.NoError(t, err, "Unexpected error")
				assert.NotNil(t, getNewsResponse.NewsFeeds, "News feed is nil")

				// Check the number of returned news feed items
				assert.Len(t, getNewsResponse.NewsFeeds, 2, "Unexpected number of news feed items")

				// Check the titles of news feed items
				expectedTitles := []string{"Item 1 Title", "Item 2 Title"}
				for i, item := range getNewsResponse.NewsFeeds {
					assert.Equal(t, expectedTitles[i], item.Title, "Unexpected news feed title")
				}

				// Check that the provider logo URL is set correctly
				assert.Equal(t, "https://www.example.com/image.jpg", getNewsResponse.NewsFeeds[0].ProviderLogoURL, "Unexpected provider logo URL")
				assert.Equal(t, "https://www.example.com/image.jpg", getNewsResponse.NewsFeeds[1].ProviderLogoURL, "Unexpected provider logo URL")
			}

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err.(ErrArgument).Err)
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
		})
	}

	// Assert that the expected method calls were made on the mock
	mockService.AssertExpectations(t)
}

func StrPointer(str string) *string {
	return &str
}
