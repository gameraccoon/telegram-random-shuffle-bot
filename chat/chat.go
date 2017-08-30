package chat

import (
	//"github.com/gameraccoon/telegram-random-shuffle-bot/database"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
)

type Chat interface {
	SendMessage(chatId int64, message string) int64
	SendDialog(chatId int64, dialog *dialog.Dialog) int64
}
