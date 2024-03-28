package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load("/Users/uday.acharya/Projects/pismo-dev/.env")
	if err != nil {
		log.Fatal("Error loading .env file with error: ", err)
	}
}
