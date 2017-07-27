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
	textFn func(*listDialogCache) string
	// nil if the variant is always active
	isActiveFn func(*listDialogCache) bool
	process func(*processing.ProcessData) bool
}

type cachedItem struct {
	id int64
	text string
}

type listDialogCache struct {
	cachedItems []cachedItem
	currentPage int
	pagesCount int
	countOnPage int
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
				isActiveFn: isTheFirstPage,
				process: addList,
			},
			variantPrototype{
				id: "first",
				textFn: getFirstItemText,
				isActiveFn: isFirstElementVisible,
				process: addList,
			},
			variantPrototype{
				id: "second",
				textFn: getSecondItemText,
				isActiveFn: isSecondElementVisible,
				process: addList,
			},
			variantPrototype{
				id: "third",
				textFn: getThirdItemText,
				isActiveFn: isThirdElementVisible,
				process: addList,
			},
			variantPrototype{
				id: "fourth",
				textFn: getFourthItemText,
				isActiveFn: isFourthElementVisible,
				process: addList,
			},
			variantPrototype{
				id: "fifth",
				textFn: getFifthItemText,
				isActiveFn: isFifthElementVisible,
				process: addList,
			},
			variantPrototype{
				id: "back",
				text: trans("back_btn"),
				isActiveFn: isNotTheFirstPage,
				process: moveBack,
			},
			variantPrototype{
				id: "fwd",
				text: trans("fwd_btn"),
				isActiveFn: isNotTheLastPage,
				process: moveForward,
			},
		},
	})
}

func getMenuText(data *processing.ProcessData) string {
	return "Choose a list"
}

func isTheFirstPage(cahce *listDialogCache) bool {
	return cahce.currentPage == 0
}

func isNotTheFirstPage(cahce *listDialogCache) bool {
	return cahce.currentPage > 0
}

func isNotTheLastPage(cahce *listDialogCache) bool {
	return cahce.currentPage + 1 < cahce.pagesCount
}

func isFirstElementVisible(cahce *listDialogCache) bool {
	return cahce.countOnPage > 0
}

func isSecondElementVisible(cahce *listDialogCache) bool {
	return cahce.countOnPage > 1
}

func isThirdElementVisible(cahce *listDialogCache) bool {
	return cahce.countOnPage > 2
}

func isFourthElementVisible(cahce *listDialogCache) bool {
	return cahce.countOnPage > 3
}

func isFifthElementVisible(cahce *listDialogCache) bool {
	return cahce.countOnPage > 4
}

func getFirstItemText(cahce *listDialogCache) string {
	index := cahce.currentPage * 4
	return cahce.cachedItems[int64(index)].text
}

func getSecondItemText(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 1
	return cahce.cachedItems[int64(index)].text
}

func getThirdItemText(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 2
	return cahce.cachedItems[int64(index)].text
}

func getFourthItemText(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 3
	return cahce.cachedItems[int64(index)].text
}

func getFifthItemText(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 4
	return cahce.cachedItems[int64(index)].text
}

func addList(data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_enter_list_name"))
	data.Static.SetUserStateTextProcessor(data.UserId, "list_name")
	return true
}

func moveForward(data *processing.ProcessData) bool {
	data.Static.SetUserStateCurrentPage(data.UserId, data.Static.GetUserStateCurrentPage(data.UserId) + 1)
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func moveBack(data *processing.ProcessData) bool {
	data.Static.SetUserStateCurrentPage(data.UserId, data.Static.GetUserStateCurrentPage(data.UserId) - 1)
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func (factory *listDialogFactory) createVariants(userId int64, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)
	cache := getCache(userId, staticData)
	
	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(cache) {
			var text string
			
			if variant.textFn != nil {
				text = variant.textFn(cache)
			} else {
				text = variant.text
			}

			variants = append(variants, dialog.Variant{
				Id:   variant.id,
				Text: text,
			})
		}
	}
	return
}

func getCache(userId int64, staticData *processing.StaticProccessStructs) (cache *listDialogCache) {

	cache = &listDialogCache{}
	
	cache.cachedItems = make([]cachedItem, 0)

	ids, names := staticData.Db.GetUserLists(userId)
	if len(ids) == len(names) {
		for index, id := range ids {
			cache.cachedItems = append(cache.cachedItems, cachedItem{
				id: id,
				text: names[index],
			})
		}
	}

	cache.currentPage = staticData.GetUserStateCurrentPage(userId)
	count := len(cache.cachedItems)
	if count > 2 {
		cache.pagesCount = (count - 2) / 4 + 1
	} else {
		cache.pagesCount = 1
	}
	
	cache.countOnPage = count - cache.currentPage * 4
	if cache.countOnPage > 5 {
		cache.countOnPage = 4
	}

	return
}

func (factory *listDialogFactory) MakeDialog(userId int64, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.text,
		Variants: factory.createVariants(userId, staticData),
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
