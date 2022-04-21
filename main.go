package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/reaganiwadha/arah/common/middleware"
	hc_client "github.com/reaganiwadha/arah/impl/hcaptcha/client"
	link_handler "github.com/reaganiwadha/arah/impl/link/http"
	"github.com/reaganiwadha/arah/impl/link/repository"
	link_uc "github.com/reaganiwadha/arah/impl/link/usecase"
	lr "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
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
	lr.Info("building server...")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.TraceId())

	lr.Info("connecting to mongo...")
	mClient, err := mongo.NewClient(options.Client().ApplyURI(config.Mongo.URI))

	if err != nil {
		lr.Fatalf("error while creating mongo mClient : %s", err)
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), time.Second*10)
	err = mClient.Connect(connectCtx)
	defer connectCancel()

	if err != nil {
		lr.Fatalf("error while connecting to mongo : %s", err)
	}

	lr.Info("connected to mongo!")

	mdb := mClient.Database(config.Mongo.Database)

	linkRepo, err := repository.NewLinkRepository(mdb)

	if err != nil {
		log.Fatalf("error while creating link repository : %s", err)
	}

	linkUc := link_uc.NewLinkUsecase(linkRepo)
	hcClient := hc_client.NewHCaptchaClient(config.Captcha.Secret)
	link_handler.ConfigureLinkHandler(linkUc, config.Server.Domain, config.Captcha.Sitekey, hcClient,
		r)

	return r
}
