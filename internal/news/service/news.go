package service

import (
	"context"
	"github.com/fir1/news/internal/news/model"
)

/*
For context, the mobile app has the following functionality:
- Load news articles from a public news feed
- Display a scrollable list of news articles
- Provide the option to filter news articles by category (such as UK News and Technology news),
where this information is available
- Show a single news article on screen, using an HTML display
- Provide the option to share news articles via email and/or social networks
- Display a thumbnail of each article in the list of articles
- Present news articles in the order in which they are published
- Allow the selection of different sources of news by category and provider
*/

type NewsInterface interface {
	GetArticle(ctx context.Context, articleURL string) (model.Article, error)
	ListNews(ctx context.Context, params ListNewsParams) (ListNewsResponse, error)
}
