package chat

import (
	//"github.com/gameraccoon/telegram-random-shuffle-bot/database"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
)

type Chat interface {
	SendMessage(chatId int64, message string)
	//SendQuestion(db *database.Database, questionId int64, usersChatIds []int64)
	SendDialog(dialog *dialog.Dialog, chatId int64)
}
