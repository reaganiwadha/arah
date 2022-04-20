package arah

import "github.com/reaganiwadha/arah/domain"

type linkUsecase struct {
}

func NewLinkUsecase() domain.LinkUsecase {
	return &linkUsecase{}
}

func (u *linkUsecase) CreateLink(slug string, link string) (res *domain.ShortenedLink, err error) {
	return
}

func (u *linkUsecase) GetLink(slug string) (res *domain.ShortenedLink, err error) {
	return
}
