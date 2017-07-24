package dialogFactory

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
)

type DialogFactory interface {
	MakeDialog(data *processing.ProcessData) *dialog.Dialog
	ProcessVariant(variantId string, data *processing.ProcessData) bool
}
