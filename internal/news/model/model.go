package model

import "time"

type NewsProvider string

const (
	NewsProviderSky   = "sky"
	NewsProviderBBC   = "bbc"
	NewsProviderOther = "other"
)

func (np NewsProvider) Valid() bool {
	switch np {
	case NewsProviderSky:
		return true
	case NewsProviderBBC:
		return true
	case NewsProviderOther:
		return true
	}
	return false
}

type NewsFeed struct {
	Title           string
	Description     string
	Link            string
	PublishDate     time.Time
	Provider        NewsProvider
	ProviderLogoURL string
}

type ByPublishDateDESC []NewsFeed

func (b ByPublishDateDESC) Len() int           { return len(b) }
func (b ByPublishDateDESC) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPublishDateDESC) Less(i, j int) bool { return b[i].PublishDate.After(b[j].PublishDate) }

type ByPublishDateASC []NewsFeed

func (b ByPublishDateASC) Len() int           { return len(b) }
func (b ByPublishDateASC) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPublishDateASC) Less(i, j int) bool { return b[i].PublishDate.Before(b[j].PublishDate) }

type Article struct {
	Title       string
	Description string
	Content     string
	Link        string
}
