package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/nicksnyder/go-i18n/i18n"
)

func MakeMenuDialogFactory(trans i18n.TranslateFunc) *DialogFactory {
	return &(DialogFactory{
		getTextFn: getMenuText,
		variants: []variantPrototype{
			variantPrototype{
				id:         "ls",
				text:       trans("list_feature_text"),
				isActiveFn: nil,
				process:    openListDialogCommand,
			},
		},
	})
}

func getMenuText(data *processing.ProcessData) string {
	return ""
}

func openListDialogCommand(data *processing.ProcessData) {
	data.Static.SetUserStateTextProcessor(data.UserId, "list_name")
}
