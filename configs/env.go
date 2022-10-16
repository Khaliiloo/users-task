package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getDBConnectionInfo() map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return map[string]string{"dburi": os.Getenv("DBURI"),
		"username": os.Getenv("USERNAME"),
		"password": os.Getenv("PASSWORD")}
}
