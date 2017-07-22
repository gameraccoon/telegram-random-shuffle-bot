package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/nicksnyder/go-i18n/i18n"
)

func MakeListsDialogFactory(trans i18n.TranslateFunc) *DialogFactory {
	return &(DialogFactory{
		getTextFn: getMenuText,
		variants: []variantPrototype{
			variantPrototype{
				id:         "add",
				text:       trans("add_list_btn"),
				isActiveFn: isFirstPage,
				process:    addList,
			},
			variantPrototype{
				id:         "first",
				text:       "item1",
				isActiveFn: nil,
				process:    addList,
			},
			variantPrototype{
				id:         "second",
				text:       "item2",
				isActiveFn: nil,
				process:    addList,
			},
			variantPrototype{
				id:         "third",
				text:       "item3",
				isActiveFn: nil,
				process:    addList,
			},
			variantPrototype{
				id:         "back",
				text:       trans("back_btn"),
				isActiveFn: nil,
				process:    addList,
			},
			variantPrototype{
				id:         "fwd",
				text:       trans("fwd_btn"),
				isActiveFn: nil,
				process:    addList,
			},
		},
	})
}

func getMenuText(data *processing.ProcessData) string {
	return "Choose a list"
}

func isFirstPage(data *processing.ProcessData) bool {
	return true
}

func addList(data *processing.ProcessData) {
	data.Static.SetUserStateTextProcessor(data.UserId, "list_name")
}
