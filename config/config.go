package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var PrivateKey string
var Publickey string
var ApiServerUrl string
var IntegrationManagerUrl string
var SecurityUrl string

// RunMode refers to run mode.
var RunMode string

func InitEnvironmentVariables() {
	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		RunMode = "DEVELOP"
	}

	if RunMode != "PRODUCTION" {
		//Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("ERROR:", err.Error())
			return
		}
	}
	log.Println("RUN MODE:", RunMode)
	PrivateKey = os.Getenv("PRIVATE_KEY")
	Publickey = os.Getenv("PUBLIC_KEY")
	ApiServerUrl = os.Getenv("API_SERVER_URL")
	IntegrationManagerUrl = os.Getenv("INTEGRATION_MANAGER_URL")
	SecurityUrl = os.Getenv("SECURITY_URL")
}
