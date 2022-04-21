package arah

import (
	"context"
	"github.com/reaganiwadha/arah/domain"
	"net/url"
	"regexp"
)

type linkUsecase struct {
	r domain.LinkRepository
}

var slugRegex *regexp.Regexp

func init() {
	slugRegex, _ = regexp.Compile("^[a-z0-9]+(?:-[a-z0-9]+)*$")
}

func (l linkUsecase) CreateLink(ctx context.Context, slug string, link string) (res *domain.ShortenedLink, err error) {
	parsedUrl, err := url.Parse(link)

	if err != nil {
		err = domain.ErrInvalidLinkFormat
		return
	}

	if !slugRegex.MatchString(slug) {
		err = domain.ErrInvalidSlugFormat
		return
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	return l.r.CreateLink(ctx, slug, parsedUrl.String())
}

func (l linkUsecase) GetLink(ctx context.Context, slug string) (res *domain.ShortenedLink, err error) {
	return l.r.GetLink(ctx, slug)
}

func NewLinkUsecase(r domain.LinkRepository) domain.LinkUsecase {
	return &linkUsecase{
		r: r,
	}
}
