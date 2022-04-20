package routing

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/reaganiwadha/arah/domain"
	"html/template"
	"net/http"
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

func (h *linkHandler) mainPage(ctx *gin.Context) {
	pageConfig := map[string]string{
		"domain":    h.domain,
		"hcsitekey": h.hcSitekey,
	}

	indexTemplate.Execute(ctx.Writer, pageConfig)
}

func (h *linkHandler) submit(ctx *gin.Context) {
	slug := ctx.PostForm("slug")
	link := ctx.PostForm("link")
	hCaptchaResponse := ctx.PostForm("h-captcha-response")

	ok, err := h.hcClient.Verify(context.TODO(), hCaptchaResponse)

	if err != nil {
		// unknown err
		return
	}

	if !ok {
		// Error page
		return
	}

	_, err = h.u.CreateLink(slug, link)

	if err != nil {
		// handle error
		return
	}

	//indexTemplate.Execute(ctx.Writer, pageConfig)
	ctx.Status(http.StatusCreated)

	//hCaptchaRespon
}
