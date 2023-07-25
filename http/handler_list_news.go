package http

import (
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache/v3"
	newsModel "github.com/fir1/news/internal/news/model"
	newsSvc "github.com/fir1/news/internal/news/service"
	"net/http"
	"time"
)

type listNewsRequest struct {
	// possible values: "sky, bbc". By default we will take news feed from all the available providers.
	// if value `providers` filled the system will try to fetch news from the given providers
	// and please don't fill anything for `news_source_url` field because you are allowed
	// to choose to get a news feed either via choosing existing providers or by giving news_source_url.
	Providers *[]string `form:"providers"`
	// one-of: DESC - latest article will be shown first in the list, ASC - oldest article will be shown first in the list
	SortByPublishDate string `form:"sort_by_publish_date" default:"DESC"`
	// possible values: "general, technology". By default we will take news feed with category `general`.
	Categories *[]string `form:"categories"`
	// if value `news_source_url` filled the system will try to fetch news from the given `url`.
	// The url must be a valid RSS url link ending with `.xml`
	// and please don't fill anything for `providers` field because you are allowed
	// to choose to get a news feed either via choosing existing providers or by giving news_source_url
	NewsSourceURL *string `form:"news_source_url"`
} // @name ListNewsRequest

type listNewsResponse struct {
	News []News `json:"news"`
} // @name ListNewsResponse

type News struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Link            string    `json:"link"`
	PublishDate     time.Time `json:"publish_date"`
	Provider        string    `json:"provider"`
	ProviderLogoURL string    `json:"provider_logo_url"`
} // @name News

// listNews example
//
//	@Summary		List news articles from a public news feed
//	@Description	 	List news articles from a public news feed
//	@Tags News
//	@ID				news-list
//	@Accept			json
//	@Produce		json
//	@Param			query-params query ListNewsRequest false "List calendar events request"
//
// @Success      200 {object}   ListNewsResponse
//
//	@Failure      400
//
// @Failure      500
// @Router			/news [get].
func (s *Service) listNews(w http.ResponseWriter, r *http.Request) {
	request := listNewsRequest{}
	err := parseQueryParamsToStruct(r, &request)
	if err != nil {
		s.respond(w, err, 0)
		return
	}

	cacheResponse, err := s.cacheClient.Get(r.RequestURI)
	switch {
	case err == nil:
		response := listNewsResponse{}
		err = json.Unmarshal(cacheResponse, &response)
		if err != nil {
			s.respond(w, err, http.StatusInternalServerError)
			return
		}
		s.respond(w, response, http.StatusOK)
		return
	case errors.Is(err, bigcache.ErrEntryNotFound):
	default:
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	var providers *[]newsModel.NewsProvider
	if request.Providers != nil {
		prs := make([]newsModel.NewsProvider, len(*request.Providers))
		for i, pr := range *request.Providers {
			prs[i] = serializeRestNewsProviderToModel(pr)
		}
		providers = &prs
	}

	newsResponse, err := s.newsService.ListNews(r.Context(), newsSvc.ListNewsParams{
		Categories:        request.Categories,
		Providers:         providers,
		SortByPublishDate: newsSvc.Sort(request.SortByPublishDate),
		NewsSourceURL:     request.NewsSourceURL,
	})
	if err != nil {
		s.respond(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := listNewsResponse{
		News: serializeNewsToRestModel(newsResponse.NewsFeeds),
	}
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	err = s.cacheClient.Set(r.RequestURI, responseBytes)
	if err != nil {
		s.respond(w, err, http.StatusInternalServerError)
		return
	}
	s.respond(w, response, http.StatusOK)
}

func serializeNewsToRestModel(feeds []newsModel.NewsFeed) []News {
	result := make([]News, len(feeds))
	for i, feed := range feeds {
		result[i] = News{
			Title:           feed.Title,
			Description:     feed.Description,
			Link:            feed.Link,
			PublishDate:     feed.PublishDate,
			Provider:        string(feed.Provider),
			ProviderLogoURL: feed.ProviderLogoURL,
		}
	}
	return result
}

func serializeRestNewsProviderToModel(np string) newsModel.NewsProvider {
	switch np {
	case "sky":
		return newsModel.NewsProviderSky
	case "bbc":
		return newsModel.NewsProviderBBC
	case "other":
		return newsModel.NewsProviderOther
	}
	return newsModel.NewsProvider(np)
}
