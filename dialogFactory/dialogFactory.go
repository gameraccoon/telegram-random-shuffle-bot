package dialogFactory

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
)

type DialogFactory interface {
	MakeDialog(id int64, staticData *processing.StaticProccessStructs) *dialog.Dialog
	ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool
}
