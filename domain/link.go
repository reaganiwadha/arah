package domain

import (
	"context"
)

type ShortenedLink struct {
	ID   string `json:"id,omitempty"`
	Slug string `json:"slug"`
	Link string `json:"link"`
}

type LinkUsecase interface {
	CreateLink(ctx context.Context, slug string, link string) (res *ShortenedLink, err error)
	GetLink(ctx context.Context, slug string) (res *ShortenedLink, err error)
}

type LinkRepository interface {
	CreateLink(ctx context.Context, slug string, link string) (res *ShortenedLink, err error)
	GetLink(ctx context.Context, slug string) (res *ShortenedLink, err error)
}
