package main

import (
	"HTOJ/templates"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/public", "./public")
	r.GET("/", func(ctx *gin.Context) {
		templates.Home().Render(ctx, ctx.Writer)
	})
	r.Run()
}
