package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

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

	id := env["DISCORD_CLIENT_ID"]
	scope := env["DISCORD_SCOPE"]

	api.GET("/login", func(c *gin.Context) {
		code, exists := c.GetQuery("code")

		scheme := "http://"
		if c.Request.TLS != nil {
			scheme = "https://"
		}

		if !exists {
			logger.Info("no code, redirecting", c.Request.URL.Scheme)

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

		c.SetCookie("authorization_token", tokens.AccessToken, tokens.ExpiresIn, "/", c.Request.Host, false, false)
		c.SetCookie("refresh_token", tokens.RefreshToken, 0, "/", c.Request.Host, false, false)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    tokens,
		})
	})
}
