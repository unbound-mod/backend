package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitializeAPI() {
	api = gin.Default()

	RegisterRoutes()

	api.Run()
}

func RegisterRoutes() {
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://unbound.rip", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterAddonRoutes(api.Group("/addons"))
	RegisterUserRoutes(api.Group("/users"))
	RegisterRepoRoutes(api.Group("/repo"))
	RegisterAuthRoutes(api.Group("/auth"))
}
