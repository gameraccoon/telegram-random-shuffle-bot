package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
)

type variantPrototype struct {
	id string
	text string
	// nil if the variant is always active
	isActiveFn func(data *processing.ProcessData) bool
	process func(data *processing.ProcessData)
}

type listDialogFactory struct {
	text string
	variants []variantPrototype
}

func MakeListsDialogFactory(trans i18n.TranslateFunc) dialogFactory.DialogFactory {
	return &(listDialogFactory{
		text: trans("test text"),
		variants: []variantPrototype{
			variantPrototype{
				id: "add",
				text: trans("add_list_btn"),
				isActiveFn: isFirstPage,
				process: addList,
			},
			variantPrototype{
				id: "first",
				text: "item1",
				isActiveFn: nil,
				process: addList,
			},
			variantPrototype{
				id: "second",
				text: "item2",
				isActiveFn: nil,
				process: addList,
			},
			variantPrototype{
				id: "third",
				text: "item3",
				isActiveFn: nil,
				process: addList,
			},
			variantPrototype{
				id: "back",
				text: trans("back_btn"),
				isActiveFn: nil,
				process: addList,
			},
			variantPrototype{
				id: "fwd",
				text: trans("fwd_btn"),
				isActiveFn: nil,
				process: addList,
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

func (factory listDialogFactory) createVariants(data *processing.ProcessData) (variants []dialog.Variant) {
	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(data) {
			if variants == nil {
				variants = make([]dialog.Variant, 0)
			}

			variants = append(variants, dialog.Variant{
				Id:   variant.id,
				Text: variant.text,
			})
		}
	}
	return
}

func (factory listDialogFactory) MakeDialog(data *processing.ProcessData) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.text,
		Variants: factory.createVariants(data),
	}
}

func (factory listDialogFactory) ProcessVariant(variantId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			variant.process(data)
			return true
		}
	}
	return false
}
