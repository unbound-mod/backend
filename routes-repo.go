package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRepoRoutes(api *gin.RouterGroup) {
	api.GET("/repo", func(c *gin.Context) {
		addons := GetAllAddons()

		c.JSON(http.StatusOK, gin.H{
			"name":        "Unbound",
			"id":          "unbound.rip",
			"description": "The official Unbound repository.",

			"iconType": "custom",
			"icon":     "debug",

			"tags": []string{
				"plugins",
				"themes",
				"icon packs",
				"fonts",
			},

			"addons": addons,
		})
	})
}
