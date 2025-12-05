package main

import (
	"HTOJ/i18n"
	"HTOJ/templates"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	bundle, err := i18n.LoadDir("i18n", "en")
	if err != nil {
		log.Fatalf("load i18n: %v", err)
	}

	r := gin.Default()
	r.Static("/public", "./public")
	r.Use(localeMiddleware(bundle))

	// API endpoint to get all translations for client-side i18n
	r.GET("/api/translations/:lang", func(ctx *gin.Context) {
		lang := ctx.Param("lang")
		if !bundle.HasLocale(lang) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "language not found"})
			return
		}
		
		translations := bundle.GetTranslations(lang)
		
		// Set cache headers - cache for 1 hour, but revalidate
		ctx.Header("Cache-Control", "public, max-age=3600, must-revalidate")
		// ETag based on translations hash for cache validation
		etag := bundle.GetETag(lang)
		ctx.Header("ETag", etag)
		
		// Check if client has cached version
		if match := ctx.GetHeader("If-None-Match"); match == etag {
			ctx.Status(http.StatusNotModified)
			return
		}
		
		ctx.JSON(http.StatusOK, translations)
	})	r.GET("/", func(ctx *gin.Context) {
		t, _, _ := fromContext(ctx, bundle)
		templates.Home(t).Render(ctx.Request.Context(), ctx.Writer)
	})

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run()
}

func localeMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := strings.ToLower(c.Query("lang"))
		if lang == "" {
			if cookie, err := c.Cookie("lang"); err == nil {
				lang = strings.ToLower(cookie)
			}
		}
		if !bundle.HasLocale(lang) {
			log.Printf("i18n: missing locale %q, fallback to %s", lang, bundle.Default())
			lang = bundle.Default()
		}
		// Persist whichever locale we resolved so the next request (without query param) stays in that language.
		c.SetCookie("lang", lang, int((30 * 24 * time.Hour).Seconds()), "/", "", false, false)
		c.Header("Content-Language", lang)
		translator := bundle.Translator(lang)
		ctx := i18n.WithContext(c.Request.Context(), translator, lang, bundle.Locales())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func fromContext(c *gin.Context, bundle *i18n.Bundle) (func(string) string, i18n.Locale, []i18n.Locale) {
	t, lang, langs := i18n.FromContext(c.Request.Context())
	if t == nil {
		if lang == "" || !bundle.HasLocale(lang) {
			lang = bundle.Default()
		}
		t = bundle.Translator(lang)
		langs = bundle.Locales()
		ctx := i18n.WithContext(context.Background(), t, lang, langs)
		c.Request = c.Request.WithContext(ctx)
	}
	return t, bundle.Current(lang), langs
}
