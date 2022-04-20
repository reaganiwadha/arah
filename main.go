package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/reaganiwadha/arah/common/middleware"
	hc_client "github.com/reaganiwadha/arah/impl/hcaptcha/client"
	link_handler "github.com/reaganiwadha/arah/impl/link/http"
	link_uc "github.com/reaganiwadha/arah/impl/link/usecase"
	lr "github.com/sirupsen/logrus"
)

func main() {
	config := loadConfig()

	r := buildServer(config)
	err := r.Run(fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		lr.Fatal("error trying to start server : %s", err)
	}
}

func buildServer(config *arahServerConfig) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.TraceId())

	linkUc := link_uc.NewLinkUsecase()
	hcClient := hc_client.NewHCaptchaClient(config.Captcha.Secret)
	link_handler.ConfigureLinkHandler(linkUc, config.Server.Domain, config.Captcha.Sitekey, hcClient,
		r)

	return r
}
