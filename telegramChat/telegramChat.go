package telegramChat

import (
	"bytes"
	"fmt"
	//"github.com/gameraccoon/telegram-random-shuffle-bot/database"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramChat struct {
	bot *tgbotapi.BotAPI
}

func MakeTelegramChat(apiToken string) (bot *TelegramChat, outErr error) {
	newBot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		outErr = err
		return
	}

	bot = &TelegramChat{
		bot: newBot,
	}

	return
}

func (telegramChat *TelegramChat) GetBot() *tgbotapi.BotAPI {
	return telegramChat.bot
}

func (telegramChat *TelegramChat) SetDebugModeEnabled(isEnabled bool) {
	telegramChat.bot.Debug = isEnabled
}

func (telegramChat *TelegramChat) GetBotUsername() string {
	return telegramChat.bot.Self.UserName
}

func (telegramChat *TelegramChat) SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = "HTML"
	telegramChat.bot.Send(msg)
}

func appendCommand(buffer *bytes.Buffer, dialogId string, variantId string, variantText string, additionalId string) {
	if additionalId == "" {
		buffer.WriteString(fmt.Sprintf("\n/%s_%s - %s", dialogId, variantId, variantText))
	} else {
		buffer.WriteString(fmt.Sprintf("\n/%s_%s_%s - %s", dialogId, variantId, additionalId, variantText))
	}
}

func (telegramChat *TelegramChat) SendDialog(chatId int64, dialog *dialog.Dialog) {
	var buffer bytes.Buffer

	buffer.WriteString(dialog.Text + "\n")

	for _, variant := range dialog.Variants {
		appendCommand(&buffer, dialog.Id, variant.Id, variant.Text, variant.AdditionalId)
	}

	telegramChat.SendMessage(chatId, buffer.String())
}
