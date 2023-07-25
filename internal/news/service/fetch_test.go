package service

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestFetchNews(t *testing.T) {
	// Creating example data for the struct
	rssData := RSS{
		Channel: Channel{
			Title:       CDATA{Text: "Sample RSS Feed"},
			Description: CDATA{Text: "This is a sample RSS feed."},
			Link:        "https://www.example.com/rss",
			Image: Image{
				URL:   "https://www.example.com/image.jpg",
				Title: "Image Title",
				Link:  "https://www.example.com",
			},
			Generator:     "My RSS Generator",
			LastBuildDate: "2023-07-25T12:00:00Z",
			Copyright:     CDATA{Text: "© 2023 Example Inc."},
			Language:      CDATA{Text: "en-us"},
			TTL:           60,
			Items: []Item{
				{
					Title:       CDATA{Text: "Item 1 Title"},
					Link:        "https://www.example.com/item1",
					Description: "This is the description of item 1.",
					GUID: GUID{
						Value:       "https://www.example.com/item1",
						IsPermaLink: "true",
					},
					PubDate: "2023-07-25T08:00:00Z",
				},
				{
					Title:       CDATA{Text: "Item 2 Title"},
					Link:        "https://www.example.com/item2",
					Description: "This is the description of item 2.",
					GUID: GUID{
						Value:       "https://www.example.com/item2",
						IsPermaLink: "true",
					},
					PubDate: "2023-07-25T10:00:00Z",
				},
			},
		},
	}

	responseBody, err := json.Marshal(&rssData)
	if err != nil {
		t.Error(err)
		return
	}

	// Create test cases using table-driven testing
	testCases := []struct {
		name         string
		statusCode   int
		retryAfter   string
		responseBody []byte
	}{
		{
			name:         "Success",
			statusCode:   http.StatusOK,
			responseBody: responseBody,
		},
		{
			name:         "RetryWithSuccess",
			statusCode:   http.StatusOK,
			retryAfter:   "2", // after 2 retries the final status code will be OK
			responseBody: responseBody,
		},
		{
			name:       "PageNotFound",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			attempts := 2 // server succeeds after 2 attempts
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// inform the client to retry after two second using standard
				// HTTP 429 status code with Retry-After header in seconds
				if tc.name == "RetryWithSuccess" && attempts > 0 {
					attempts--
					w.Header().Set("Retry-After", tc.retryAfter)
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("Server limit reached"))
					return
				}

				w.Write(tc.responseBody)

				w.WriteHeader(tc.statusCode)
			}))
			defer ts.Close()

			var body []byte
			err = retry.Do(
				func() error {
					resp, err := http.Get(ts.URL)
					if err != nil {
						return err
					}

					defer resp.Body.Close()
					body, err = io.ReadAll(resp.Body)
					if resp.StatusCode != http.StatusOK {
						err = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
						if resp.StatusCode == http.StatusTooManyRequests {
							// check Retry-After header if it contains seconds to wait for the next retry
							if retryAfter, e := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 32); e == nil {
								// the server returns 0 to inform that the operation cannot be retried
								if retryAfter <= 0 {
									return retry.Unrecoverable(err)
								}
								return &RetriableError{
									Err:        err,
									RetryAfter: time.Duration(retryAfter) * time.Second,
								}
							}
						}
					}
					return nil
				},
				retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
					fmt.Println("Server fails with: " + err.Error())
					if retriable, ok := err.(*RetriableError); ok {
						fmt.Printf("Client follows server recommendation to retry after %v\n", retriable.RetryAfter)
						return retriable.RetryAfter
					}
					// apply a default exponential back off strategy
					return retry.BackOffDelay(n, err, config)
				}),
			)

			fmt.Println("Server responds with: " + string(body))

			if tc.responseBody != nil {
				assert.Equal(t, responseBody, body)
			}
		})
	}
}
