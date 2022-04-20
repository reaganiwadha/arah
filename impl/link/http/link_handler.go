package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/reaganiwadha/arah/domain"
	"html/template"
)

//var

func ConfigureLinkRouting(u domain.LinkUsecase, r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		indexTemplate := template.Must(template.ParseFiles("templates/index.gohtml"))
		indexTemplate.Execute(ctx.Writer, nil)
	})
}
