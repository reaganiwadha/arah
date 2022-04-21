package routing

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/reaganiwadha/arah/domain"
	"html/template"
	"net/http"
	"time"
)

var indexTemplate = template.Must(template.ParseFiles("templates/index.gohtml"))

type linkHandler struct {
	u         domain.LinkUsecase
	domain    string
	hcSitekey string
	hcClient  domain.HCaptchaClient
}

func ConfigureLinkHandler(u domain.LinkUsecase, domain string, hcSitekey string, hCaptchaClient domain.HCaptchaClient, r *gin.Engine) {
	handler := linkHandler{
		u:         u,
		domain:    domain,
		hcSitekey: hcSitekey,
		hcClient:  hCaptchaClient,
	}

	r.GET("/", handler.mainPage)
	r.POST("/submit", handler.submit)
}

type pageTemplateData struct {
	Domain        string
	HcSitekey     string
	ErrorString   string
	ShortenedLink string
	OriginalLink  string
}

func (h *linkHandler) createPageTemplateData() pageTemplateData {
	return pageTemplateData{
		Domain:    h.domain,
		HcSitekey: h.hcSitekey,
	}
}

func (h *linkHandler) mainPage(ctx *gin.Context) {
	indexTemplate.Execute(ctx.Writer, h.createPageTemplateData())
}

func (h *linkHandler) submit(ctx *gin.Context) {
	slug := ctx.PostForm("slug")
	link := ctx.PostForm("link")
	hCaptchaResponse := ctx.PostForm("h-captcha-response")

	verifyCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ok, err := h.hcClient.Verify(verifyCtx, hCaptchaResponse)
	defer cancel()

	pageConfig := h.createPageTemplateData()

	if err != nil || !ok {
		if err != nil {
			pageConfig.ErrorString = "Unknown Error"
			ctx.Status(http.StatusInternalServerError)
		} else if !ok {
			pageConfig.ErrorString = "Captcha Failed"
			ctx.Status(http.StatusBadRequest)
		}

		indexTemplate.Execute(ctx.Writer, pageConfig)
		return
	}

	clCtx, clCancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err = h.u.CreateLink(clCtx, slug, link)
	defer clCancel()

	if err != nil {
		if errors.Is(err, domain.ErrDataExists) {
			pageConfig.ErrorString = fmt.Sprintf("The short for `%s` already exists 😔", slug)
			ctx.Status(http.StatusBadRequest)
		} else if errors.Is(err, domain.ErrInvalidSlugFormat) {
			pageConfig.ErrorString = "Invalid slug format"
			ctx.Status(http.StatusBadRequest)
		} else if errors.Is(err, domain.ErrInvalidLinkFormat) {
			pageConfig.ErrorString = "Invalid link"
			ctx.Status(http.StatusBadRequest)
		} else {
			pageConfig.ErrorString = "Unknown Error"
			ctx.Status(http.StatusInternalServerError)
		}

		indexTemplate.Execute(ctx.Writer, pageConfig)
		return
	}

	pageConfig.OriginalLink = link
	pageConfig.ShortenedLink = fmt.Sprintf("https://%s/%s", h.domain, slug)

	indexTemplate.Execute(ctx.Writer, pageConfig)
	ctx.Status(http.StatusCreated)
}
