package arah

import (
	"context"
	"github.com/reaganiwadha/arah/domain"
)

type linkUsecase struct {
	r domain.LinkRepository
}

func (l linkUsecase) CreateLink(ctx context.Context, slug string, link string) (res *domain.ShortenedLink, err error) {
	return l.r.CreateLink(ctx, slug, link)
}

func (l linkUsecase) GetLink(ctx context.Context, slug string) (res *domain.ShortenedLink, err error) {
	return
}

func NewLinkUsecase(r domain.LinkRepository) domain.LinkUsecase {
	return &linkUsecase{
		r: r,
	}
}
