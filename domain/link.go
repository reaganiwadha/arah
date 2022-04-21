package domain

import (
	"context"
	"errors"
)

type ShortenedLink struct {
	ID   string `json:"id,omitempty"`
	Slug string `json:"slug"`
	Link string `json:"link"`
}

var (
	ErrInvalidLinkFormat = errors.New("invalid link format")
	ErrInvalidSlugFormat = errors.New("invalid slug format")
)

type LinkUsecase interface {
	CreateLink(ctx context.Context, slug string, link string) (res *ShortenedLink, err error)
	GetLink(ctx context.Context, slug string) (res *ShortenedLink, err error)
}

type LinkRepository interface {
	CreateLink(ctx context.Context, slug string, link string) (res *ShortenedLink, err error)
	GetLink(ctx context.Context, slug string) (res *ShortenedLink, err error)
}
