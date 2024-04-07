package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	Structures "unbound.rip/backend/structures"
)

func RegisterAuthRoutes(api *gin.RouterGroup) {
	id := os.Getenv("DISCORD_CLIENT_ID")
	scope := os.Getenv("DISCORD_SCOPE")
	domain := os.Getenv("DOMAIN")

	if isDevelopment {
		domain = "localhost"
	}

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

		redirect := scheme + c.Request.Host + "/auth/login"
		tokens, err := GetAuthorizationFromCode(code, redirect)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprint(err),
			})

			return
		}

		user, err := GetDiscordUserFromAuth(tokens.AccessToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprint(err),
			})

			return
		}

		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "User returned by Discord was nil.",
			})

			return
		}

		// TODO: Figure out how to upsert the db record for an existing user.

		// payload := &Structures.User{
		// 	Discord: *user,
		// 	Tokens: Structures.UserTokens{
		// 		AccessToken:  tokens.AccessToken,
		// 		RefreshToken: tokens.RefreshToken,
		// 	},
		// }

		// db.Clauses(clause.OnConflict{
		// 	Columns: []clause.Column{{Name: "discord_id"}},
		// 	DoUpdates: clause.AssignmentColumns([]string{
		// 		"discord_avatar",
		// 		"discord_username",
		// 		"discord_display_name",
		// 		"tokens_access_token",
		// 		"tokens_refresh_token",
		// 	}),
		// }).Create(payload)

		c.SetCookie("token", tokens.AccessToken, 0, "/", domain, false, false)
		c.SetCookie("refresh_token", tokens.RefreshToken, 0, "/", domain, false, false)

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

		redirect := scheme + c.Request.Host + "/auth/login"
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

func GetAuth(c *gin.Context) (*Structures.User, error) {
	header := c.GetHeader("Authorization")

	if header == "" {
		return nil, Structures.ErrCredentialsMissing
	}

	user := GetUserFromAuth(header)

	return user, nil
}

func GetUserFromAuth(auth string) *Structures.User {
	var user *Structures.User

	db.Find(&user)

	return user
}
