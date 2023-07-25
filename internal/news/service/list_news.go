package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/fir1/news/internal/news/model"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type ListNewsParams struct {
	Categories        *[]string
	Providers         *[]model.NewsProvider
	NewsSourceURL     *string
	SortByPublishDate Sort
}

type ListNewsResponse struct {
	NewsFeeds []model.NewsFeed
}

func (s Service) ListNews(ctx context.Context, params ListNewsParams) (ListNewsResponse, error) {
	if params.Categories == nil {
		params.Categories = &[]string{"general"}
	}

	if params.SortByPublishDate != "" && !params.SortByPublishDate.Valid() {
		return ListNewsResponse{}, ErrArgument{Err: errors.New("please provide a valid sort by publish date ASC or DESC")}
	}

	if params.Providers != nil && params.NewsSourceURL != nil {
		return ListNewsResponse{}, ErrArgument{Err: errors.New("please provide one of value for providers or news_source_url can not proceed both")}
	}

	newsFeedCh := make(chan []model.NewsFeed)
	errorCh := make(chan error)
	var wg sync.WaitGroup

	// in case both providers and new_source_url not provided by client, we will take all available news_providers by default
	if params.Providers == nil && params.NewsSourceURL == nil {
		params.Providers = &[]model.NewsProvider{model.NewsProviderBBC, model.NewsProviderSky}
	}

	if params.Providers != nil {
		for _, provider := range *params.Providers {
			if !provider.Valid() {
				return ListNewsResponse{}, ErrArgument{Err: fmt.Errorf("provider: %s is invalid must be `sky`, `bbc`", provider)}
			}

			for _, category := range *params.Categories {
				if category != "general" && category != "technology" {
					return ListNewsResponse{}, ErrArgument{Err: fmt.Errorf("category: %s is invalid must be `general`, `technology`", category)}
				}
				wg.Add(1)

				feedURL, found := feedURLs[provider][category]
				if !found {
					feedURL = feedURLs[provider]["general"]
				}

				go func(wg *sync.WaitGroup, feedURL string, provider model.NewsProvider) {
					defer wg.Done()
					newsFeed, err := s.getProviderNewsFeed(ctx, feedURL, provider)
					if err != nil {
						errorCh <- err
						return
					}
					newsFeedCh <- newsFeed
				}(&wg, feedURL, provider)
			}
		}
	}

	if params.NewsSourceURL != nil {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			newsFeed, err := s.getProviderNewsFeed(ctx, *params.NewsSourceURL, model.NewsProviderOther)
			if err != nil {
				errorCh <- err
				return
			}
			newsFeedCh <- newsFeed
		}(&wg)
	}

	// Wait for all goroutines to finish and close the result channel
	go func() {
		wg.Wait()
		close(newsFeedCh)
	}()

	var combinedResult ListNewsResponse
	run := true
	for run {
		select {
		case result, ok := <-newsFeedCh:
			if ok {
				combinedResult.NewsFeeds = append(combinedResult.NewsFeeds, result...)
			} else {
				// resultCh is closed, we are done processing results
				run = false
			}
		case err, ok := <-errorCh:
			if ok {
				return ListNewsResponse{}, err
			}
			// errorCh is closed, we are done processing errors
			run = false
		}
	}

	if params.SortByPublishDate == SortASC {
		sort.Sort(model.ByPublishDateASC(combinedResult.NewsFeeds))
	} else {
		// Sort the News by PublishDate (latest first) by DEFAULT
		sort.Sort(model.ByPublishDateDESC(combinedResult.NewsFeeds))
	}
	return combinedResult, nil
}

var feedURLs = map[model.NewsProvider]map[string]string{
	model.NewsProviderBBC: {
		"general":    "http://feeds.bbci.co.uk/news/uk/rss.xml",
		"technology": "http://feeds.bbci.co.uk/news/technology/rss.xml",
	},
	model.NewsProviderSky: {
		"general":    "http://feeds.skynews.com/feeds/rss/uk.xml",
		"technology": "http://feeds.skynews.com/feeds/rss/technology.xml",
	},
}

func (s Service) getProviderNewsFeed(ctx context.Context, feedURL string, provider model.NewsProvider) ([]model.NewsFeed, error) {
	if ok := isValidURL(feedURL); !ok {
		return nil, ErrArgument{Err: fmt.Errorf("url: %s not valid, please provide a valid url ending with `.xml`", feedURL)}
	}

	var response []model.NewsFeed
	var feeds RSS
	err := retry.Do(
		func() error {
			var err error
			feeds, err = s.NewsFetcher.fetchNews(ctx, feedURL)
			return err
		},

		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if retriable, ok := err.(*RetriableError); ok {
				return retriable.RetryAfter
			}
			// apply a default exponential back off strategy
			return retry.BackOffDelay(n, err, config)
		}),
	)
	if err != nil {
		return nil, err
	}

	for _, item := range feeds.Channel.Items {
		pubDate, err := parseTimeFromString(item.PubDate)
		if err != nil {
			return nil, err
		}

		newsFeed := model.NewsFeed{
			Title:           item.Title.Text,
			Description:     item.Description,
			Link:            item.Link,
			PublishDate:     pubDate,
			Provider:        provider,
			ProviderLogoURL: feeds.Channel.Image.URL,
		}
		response = append(response, newsFeed)
	}

	return response, nil
}

func isValidURL(input string) bool {
	// Parse the input string as a URL
	u, err := url.Parse(input)
	if err != nil {
		return false
	}

	// Check if the URL has a valid scheme (http, https, etc.)
	if u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}

	// Check if the URL ends with ".xml"
	if !strings.HasSuffix(u.Path, ".xml") {
		return false
	}

	return true
}

func parseTimeFromString(input string) (time.Time, error) {
	layouts := []string{
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02",                     // Format: yyyy-MM-dd
		"02-01-2006",                     // Format: dd-MM-yyyy
		"2006-01-02 15:04:05",            // Format: yyyy-MM-dd HH:mm:ss
		"02-01-2006 15:04:05",            // Format: dd-MM-yyyy HH:mm:ss
		"Jan 02, 2006",                   // Format: MMM dd, yyyy
		"January 02, 2006",               // Format: Month dd, yyyy
		"Mon, Jan 02, 2006",              // Format: Mon, MMM dd, yyyy
		"Mon, January 02, 2006",          // Format: Mon, Month dd, yyyy
		"Jan 02, 2006 15:04:05",          // Format: MMM dd, yyyy HH:mm:ss
		"January 02, 2006 15:04:05",      // Format: Month dd, yyyy HH:mm:ss
		"Mon, Jan 02, 2006 15:04:05",     // Format: Mon, MMM dd, yyyy HH:mm:ss
		"Mon, January 02, 2006 15:04:05", // Format: Mon, Month dd, yyyy HH:mm:ss
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, input)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time from input string: %s", input)
}
