package service

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"time"
)

type NewsFetcher interface {
	fetchNews(ctx context.Context, feedURL string) (RSS, error)
}

type RealFetcherService struct {
}

func NewRealFetcherService() NewsFetcher {
	return RealFetcherService{}
}

// FetchNewsFeeds fetches news articles from the given feed URL and returns a slice of NewsFeed objects.
func (s RealFetcherService) fetchNews(ctx context.Context, feedURL string) (RSS, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return RSS{}, err
	}

	client := http.Client{
		Timeout: 1 * time.Minute,
	}

	resp, err := client.Do(request)
	if err != nil {
		return RSS{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == http.StatusTooManyRequests {
			// check Retry-After header if it contains seconds to wait for the next retry
			retryAfter, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 32)
			if err != nil {
				return RSS{}, err
			}

			// the server returns 0 to inform that the operation cannot be retried
			if retryAfter <= 0 {
				return RSS{}, err
			}

			return RSS{}, &RetriableError{
				Err:        err,
				RetryAfter: time.Duration(retryAfter) * time.Second,
			}
		}
	}

	var result RSS
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return RSS{}, err
	}

	return result, nil
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         CDATA  `xml:"title"`
	Description   CDATA  `xml:"description"`
	Link          string `xml:"link"`
	Image         Image  `xml:"image"`
	Generator     string `xml:"generator"`
	LastBuildDate string `xml:"lastBuildDate"`
	Copyright     CDATA  `xml:"copyright"`
	Language      CDATA  `xml:"language"`
	TTL           int    `xml:"ttl"`
	Items         []Item `xml:"item"`
}

type CDATA struct {
	Text string `xml:",cdata"`
}

type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Item struct {
	Title       CDATA  `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        GUID   `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	// Add more fields here if needed for the <item> element
}

type GUID struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}
