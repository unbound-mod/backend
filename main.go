package main

import (
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    false,
	ReportTimestamp: true,
	TimeFormat:      time.TimeOnly,
	Level:           log.DebugLevel,
	Prefix:          "Backend",
})

var (
	db            *gorm.DB
	api           *gin.Engine
	isDevelopment bool
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	request       = &http.Client{}
)

func main() {
	if os.Getenv("MODE") == "PRODUCTION" {
		logger.Info("Running in PRODUCTION mode.")
		gin.SetMode(gin.ReleaseMode)
	} else {
		err := godotenv.Load()
		logger.Info("Running in DEVELOPMENT mode.")
		isDevelopment = true

		if err != nil {
			logger.Fatalf("Failed to load .env file (Does it exist?): %v", err)
			return
		}
	}
	InitializeORM()
	InitializeAPI()
}
