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
	isActiveFn func(*listDialogFactory, *processing.ProcessData) bool
	process func(*processing.ProcessData) bool
}

type listDialogFactory struct {
	text string
	variants []variantPrototype
	cachedItems map[int64]string
	currentPage int
	pagesCount int
	isInited bool
}

func MakeListsDialogFactory(trans i18n.TranslateFunc) dialogFactory.DialogFactory {
	return &(listDialogFactory{
		text: trans("test text"),
		variants: []variantPrototype{
			variantPrototype{
				id: "add",
				text: trans("add_list_btn"),
				isActiveFn: isTheFirstPage,
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
				isActiveFn: isNotTheFirstPage,
				process: addList,
			},
			variantPrototype{
				id: "fwd",
				text: trans("fwd_btn"),
				isActiveFn: isNotTheLastPage,
				process: addList,
			},
		},
	})
}

func getMenuText(data *processing.ProcessData) string {
	return "Choose a list"
}

func isTheFirstPage(factory *listDialogFactory, data *processing.ProcessData) bool {
	factory.cacheItems(data)
	return factory.currentPage == 0
}

func isNotTheFirstPage(factory *listDialogFactory, data *processing.ProcessData) bool {
	factory.cacheItems(data)
	return factory.currentPage > 0
}

func isNotTheLastPage(factory *listDialogFactory, data *processing.ProcessData) bool {
	factory.cacheItems(data)
	return factory.currentPage + 1 < factory.pagesCount
}

func addList(data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_enter_list_name"))
	data.Static.SetUserStateTextProcessor(data.UserId, "list_name")
	return true
}

func (factory *listDialogFactory) createVariants(data *processing.ProcessData) (variants []dialog.Variant) {
	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(factory, data) {
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

func (factory *listDialogFactory) cacheItems(data *processing.ProcessData) {
	if (factory.isInited) {
		return
	}

	factory.cachedItems = make(map[int64]string)

	ids, names := data.Static.Db.GetUserLists(data.UserId)
	if len(ids) == len(names) {
		for index, id := range ids {
			factory.cachedItems[id] = names[index]
		}
	}

	factory.currentPage = data.Static.GetUserStateCurrentPage(data.UserId)
	count := len(factory.cachedItems)
	if count > 2 {
		factory.pagesCount = (count - 2) / 4 + 1
	} else {
		factory.pagesCount = 1
	}

	factory.isInited = true
}

func (factory *listDialogFactory) MakeDialog(data *processing.ProcessData) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.text,
		Variants: factory.createVariants(data),
	}
}

func (factory *listDialogFactory) ProcessVariant(variantId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			variant.process(data)
			return true
		}
	}
	return false
}
