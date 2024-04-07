package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Structures "unbound.rip/backend/structures"
)

func RegisterAddonRoutes(api *gin.RouterGroup) {
	api.GET("/:type", func(c *gin.Context) {
		addonType, exists := c.Params.Get("type")

		if !exists || (addonType != "themes" && addonType != "plugins") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "The addon type you requested is not valid. Valid types are: themes, plugins",
			})

			return
		}

		addons := GetAddons(addonType)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    addons,
		})
	})
}

func GetAllAddons() []Structures.Addon {
	var addons []Structures.Addon

	db.Find(&addons)

	return addons
}

func GetAddons(addonType string) []Structures.Addon {
	var addons []Structures.Addon

	db.Where("type = ?", addonType).Find(&addons)

	return addons
}
