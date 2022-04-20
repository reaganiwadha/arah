package middleware

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	lr "github.com/sirupsen/logrus"
)

func TraceId() gin.HandlerFunc {
	node, err := snowflake.NewNode(1)

	if err != nil {
		lr.Fatal(err)
	}

	return func(c *gin.Context) {
		traceId := node.Generate()

		c.Set("trace-id", traceId.String())
		c.Header("X-Req-Id", traceId.String())
		c.Next()
	}
}
