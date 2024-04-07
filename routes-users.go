package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	Structures "unbound.rip/backend/structures"
)

func RegisterUserRoutes(api *gin.RouterGroup) {
	api.GET("/", func(c *gin.Context) {
		if !isDevelopment {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   fmt.Sprint(Structures.ErrDeveloperOnly),
			})

			return
		}

		var result []Structures.User

		db.Find(&result)

		_, err := json.Marshal(result)

		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"success": false,
				"error":   fmt.Sprint(err),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	api.GET("/developers", func(c *gin.Context) {
		var result []Structures.User

		db.Where("developer = true").Find(&result)

		_, err := json.Marshal(result)

		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"success": false,
				"error":   fmt.Sprint(err),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	api.GET("/:idOrUsername", func(c *gin.Context) {
		param, _ := c.Params.Get("idOrUsername")

		var result Structures.User

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
}
