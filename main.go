package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"aiart-bot/bot"
	"aiart-bot/config"
)

func main() {
	config.LoadEnv()

	botAPI, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	log.Printf("Bot started: %s", botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Напиши описание — и я сгенерирую изображение")
				botAPI.Send(msg)
			} else {
				bot.HandleMessage(botAPI, update.Message)
			}
		}
	}
}
