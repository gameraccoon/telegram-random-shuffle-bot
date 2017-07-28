package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
)

type listItemVariantPrototype struct {
	id string
	text string
	process func(string, *processing.ProcessData) bool
}

type listItemDialogFactory struct {
	variants []listItemVariantPrototype
}

func MakeListItemDialogFactory(trans i18n.TranslateFunc) dialogFactory.DialogFactory {
	return &(listItemDialogFactory{
		variants: []listItemVariantPrototype{
			listItemVariantPrototype{
				id: "random",
				text: trans("get_random"),
				process: getRandom,
			},
			listItemVariantPrototype{
				id: "shuffled",
				text: trans("get_shuffled"),
				process: getShuffled,
			},
			listItemVariantPrototype{
				id: "additems",
				text: trans("add_btn"),
				process: addItems,
			},
			listItemVariantPrototype{
				id: "back",
				text: trans("back_btn"),
				process: backToLists,
			},
		},
	})
}

func getRandom(additionalId string, data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("test1"))
	return true
}

func getShuffled(additionalId string, data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("test2"))
	return true
}

func addItems(additionalId string, data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("test3"))
	return true
}

func backToLists(additionalId string, data *processing.ProcessData) bool {
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func (factory *listItemDialogFactory) getListItemDialogText(listId int64, staticData *processing.StaticProccessStructs) string {
	return staticData.Db.GetListName(listId)
}

func (factory *listItemDialogFactory) createVariants(listId int64, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)
	
	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:   variant.id,
			Text: variant.text,
		})
	}
	return
}

func (factory *listItemDialogFactory) MakeDialog(listId int64, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.getListItemDialogText(listId, staticData),
		Variants: factory.createVariants(listId, staticData),
	}
}

func (factory *listItemDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			variant.process(additionalId, data)
			return true
		}
	}
	return false
}
