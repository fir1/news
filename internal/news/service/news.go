package service

import (
	"context"
	"github.com/fir1/news/internal/news/model"
)

type NewsInterface interface {
	GetArticle(ctx context.Context, articleURL string) (model.Article, error)
	ListNews(ctx context.Context, params ListNewsParams) (ListNewsResponse, error)
}
