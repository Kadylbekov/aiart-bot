package bot

import (
	"log"

	"aiart-bot/ai"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userPrompt := message.Text
	log.Printf("Получен запрос: %s от %s", userPrompt, message.From.UserName)

	bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadPhoto))

	imgBytes, err := ai.GenerateImage(userPrompt)
	if err != nil {
		log.Printf("Ошибка генерации изображения: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Не удалось сгенерировать изображение 😔")
		bot.Send(msg)
		return
	}

	photo := tgbotapi.FileBytes{
		Name:  "anime.png",
		Bytes: imgBytes,
	}
	msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	msg.Caption = "Вот что получилось!"
	bot.Send(msg)
}
