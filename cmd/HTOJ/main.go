package main

import (
	"HTOJ/templates"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	r.Static("/public", "./public")
	// Define a simple GET endpoint
	r.GET("/ping", func(ctx *gin.Context) {
		templates.Hello("World").Render(ctx, ctx.Writer)
	})
	r.GET("/dog", func(ctx *gin.Context) {
		res, err := http.Get("https://dog.ceo/api/breeds/image/random")
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("this is bad"))
			return
		}
		defer res.Body.Close()

		dogRes := make(map[string]string)
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&dogRes)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("this is bad"))
			return
		}
		templates.Image(dogRes["message"]).Render(ctx, ctx.Writer)
	})
	r.GET("/", func(ctx *gin.Context) {
		templates.Base().Render(ctx, ctx.Writer)
	})
	r.GET("/problems", func(ctx *gin.Context) {
		templates.Home().Render(ctx, ctx.Writer)
	})
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
