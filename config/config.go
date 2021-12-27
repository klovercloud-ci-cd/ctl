package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)
var PrivateKey string
var Publickey string
var ApiServerUrl string
func InitEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error(),", reading from env")

	}
	PrivateKey =os.Getenv("PRIVATE_KEY")
	Publickey=os.Getenv("PUBLIC_KEY")
	ApiServerUrl=os.Getenv("API_SERVER_URL")

}