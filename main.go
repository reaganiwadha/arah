package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	link_handler "github.com/reaganiwadha/arah/impl/link/http"
	link_uc "github.com/reaganiwadha/arah/impl/link/usecase"
	lr "github.com/sirupsen/logrus"
)

func main() {
	config := loadConfig()

	r := buildServer()
	err := r.Run(fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		lr.Fatal("error trying to start server : %s", err)
	}
}

func buildServer() *gin.Engine {
	r := gin.Default()

	linkUc := link_uc.NewLinkUsecase()
	link_handler.ConfigureLinkRouting(linkUc, r)

	return r
}
