package http

import (
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/fir1/news/internal/news/model"
	"html/template"
	"net/http"
)

type getArticleRequest struct {
	URL string `form:"url"`
} // @name GetArticleRequest

// getArticle example
//
//	@Summary		Get article, it shows a single news article on screen, using an HTML display
//	@Description	 	Get article, it shows a single news article on screen, using an HTML display
//	@Tags News
//	@ID				article-get
//	@Accept			json
//	@Produce		json
//	@Param			query-params query GetArticleRequest false "Get article query params"
//
// @Success      200
//
//	@Failure      400
//
// @Failure      500
// @Router			/article [get].
func (s *Service) getArticle(w http.ResponseWriter, r *http.Request) {
	request := getArticleRequest{}
	err := parseQueryParamsToStruct(r, &request)
	if err != nil {
		s.respond(w, err, 0)
		return
	}

	article := model.Article{}
	cacheResponse, err := s.cacheClient.Get(r.RequestURI)
	switch {
	case err == nil:
		err = json.Unmarshal(cacheResponse, &article)
		if err != nil {
			s.respond(w, err, http.StatusInternalServerError)
			return
		}
	case errors.Is(err, bigcache.ErrEntryNotFound):
		article, err = s.newsService.GetArticle(r.Context(), request.URL)
		if err != nil {
			s.respond(w, err.Error(), http.StatusBadRequest)
			return
		}

		responseBytes, err := json.Marshal(&article)
		if err != nil {
			s.respond(w, err, http.StatusInternalServerError)
			return
		}

		err = s.cacheClient.Set(r.RequestURI, responseBytes)
		if err != nil {
			s.respond(w, err, http.StatusInternalServerError)
			return
		}
	default:
		s.respond(w, err, http.StatusInternalServerError)
		return
	}

	// HTML template
	tmpl := `<html>
		<head>
			<title>{{ .Title }}</title>
		</head>
		<body>
			<h1>{{ .Title }}</h1>
			<p>{{ .Description }}</p>
			<hr>
			{{ .Content | safe }}
		</body>
		</html>`

	// Parse the HTML template
	t := template.Must(template.New("newsArticle").Funcs(template.FuncMap{
		"safe": func(text string) template.HTML {
			return template.HTML(text)
		},
	}).Parse(tmpl))

	// Execute the template with the article data and write the result to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, article)
	if err != nil {
		http.Error(w, "Failed to render the template", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
