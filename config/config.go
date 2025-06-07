package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	TelegramToken string
	NovitaAPIKey  string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TelegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	NovitaAPIKey = os.Getenv("NOVITA_API_KEY")

	if TelegramToken == "" || NovitaAPIKey == "" {
		log.Fatal("Missing required tokens")
	}
}
