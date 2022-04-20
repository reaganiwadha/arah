package domain

type ShortenedLink struct {
}

type LinkUsecase interface {
	CreateLink(slug string, link string) (res *ShortenedLink, err error)
	GetLink(slug string) (res *ShortenedLink, err error)
}
