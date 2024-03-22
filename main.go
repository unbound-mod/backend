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
	env           map[string]string
	isDevelopment bool
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	request       = &http.Client{}
)

func main() {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		logger.Fatalf("Failed to load .env file (Does it exist?): %v", err)
		return
	}

	// Go does not allow assigning to global variables while also creating new variables
	env = envFile

	if env["MODE"] == "DEVELOPMENT" {
		logger.Info("Running in DEVELOPMENT mode.")
		isDevelopment = true
	} else {
		logger.Info("Running in PRODUCTION mode.")
		gin.SetMode(gin.ReleaseMode)
	}

	InitializeORM()
	InitializeAPI()
}
