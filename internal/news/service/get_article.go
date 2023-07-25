package service

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fir1/news/internal/news/model"
	"net/http"
	"net/url"
	"strings"
)

func (s Service) GetArticle(ctx context.Context, articleURL string) (model.Article, error) {
	// Validate the URL
	_, err := url.Parse(articleURL)
	if err != nil {
		return model.Article{}, ErrArgument{Err: fmt.Errorf("invalid URL: %v", err)}
	}

	resp, err := http.Get(articleURL)
	if err != nil {
		return model.Article{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Article{}, fmt.Errorf("failed to fetch article status code %d", resp.StatusCode)
	}

	// Parse the HTML content of the article
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return model.Article{}, err
	}

	// Extract the title, description, and content of the article
	var title, description, content string
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		title = strings.TrimSpace(s.Text())
	})

	doc.Find("meta[name=description]").Each(func(i int, s *goquery.Selection) {
		description, _ = s.Attr("content")
	})

	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		content = strings.TrimSpace(s.Text())
	})

	return model.Article{
		Title:       title,
		Description: description,
		Link:        articleURL,
		Content:     content,
	}, nil
}
