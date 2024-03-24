package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeAPI() {
	api = gin.Default()

	RegisterRoutes()

	api.Run()
}

func RegisterRoutes() {
	api.GET("/users/:idOrUsername", func(c *gin.Context) {
		param, _ := c.Params.Get("idOrUsername")

		var result User

		query := db.First(&result, param)

		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "User not found.",
			})

			return
		}

		if query.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Database lookup failed.",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	id := os.Getenv("DISCORD_CLIENT_ID")
	scope := os.Getenv("DISCORD_SCOPE")
	domain := os.Getenv("DOMAIN")

	if isDevelopment {
		domain = "localhost"
	}

	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://unbound.rip", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api.GET("/login", func(c *gin.Context) {
		code, exists := c.GetQuery("code")

		scheme := "http://"
		if c.Request.TLS != nil {
			scheme = "https://"
		}

		if !exists {
			redirect := url.QueryEscape(scheme + c.Request.Host + c.Request.URL.Path)
			scopes := url.QueryEscape(scope)

			url := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%v&scope=%v&permissions=0&response_type=code&redirect_uri=%v", id, scopes, redirect)
			c.Redirect(http.StatusTemporaryRedirect, url)

			return
		}

		redirect := scheme + c.Request.Host + "/login"
		tokens, err := GetAuthorizationFromCode(code, redirect)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprint(err),
			})

			return
		}

		c.SetCookie("authorization_token", tokens.AccessToken, tokens.ExpiresIn, "/", domain, false, false)
		c.SetCookie("refresh_token", tokens.RefreshToken, 34560000, "/", domain, false, false)

		if isDevelopment {
			c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000")
		} else {
			c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("https://%v", domain))
		}
	})

	api.POST("/revoke", func(c *gin.Context) {
		token, exists := c.GetQuery("token")

		scheme := "http://"
		if c.Request.TLS != nil {
			scheme = "https://"
		}

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "URL param \"token\" is not present.",
			})

			return
		}

		redirect := scheme + c.Request.Host + "/login"
		err := RevokeAuthorization(token, redirect)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprint(err),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}
