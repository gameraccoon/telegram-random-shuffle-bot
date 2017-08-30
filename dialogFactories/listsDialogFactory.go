package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
	"strconv"
)

type listDialogVariantPrototype struct {
	id string
	additionalIdFn func(*listDialogCache) string
	text string
	textFn func(*listDialogCache) string
	// nil if the variant is always active
	isActiveFn func(*listDialogCache) bool
	process func(string, *processing.ProcessData) bool
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
	variants []listDialogVariantPrototype
}

func MakeListsDialogFactory(trans i18n.TranslateFunc) dialogFactory.DialogFactory {
	return &(listDialogFactory{
		text: trans("choose_list"),
		variants: []listDialogVariantPrototype{
			listDialogVariantPrototype{
				id: "add",
				text: trans("add_list_btn"),
				isActiveFn: isTheFirstPage,
				process: addList,
			},
			listDialogVariantPrototype{
				id: "it1",
				additionalIdFn: getFirstItemId,
				textFn: getFirstItemText,
				isActiveFn: isFirstElementVisible,
				process: openListItem,
			},
			listDialogVariantPrototype{
				id: "it2",
				additionalIdFn: getSecondItemId,
				textFn: getSecondItemText,
				isActiveFn: isSecondElementVisible,
				process: openListItem,
			},
			listDialogVariantPrototype{
				id: "it3",
				additionalIdFn: getThirdItemId,
				textFn: getThirdItemText,
				isActiveFn: isThirdElementVisible,
				process: openListItem,
			},
			listDialogVariantPrototype{
				id: "it4",
				additionalIdFn: getFourthItemId,
				textFn: getFourthItemText,
				isActiveFn: isFourthElementVisible,
				process: openListItem,
			},
			listDialogVariantPrototype{
				id: "it5",
				additionalIdFn: getFifthItemId,
				textFn: getFifthItemText,
				isActiveFn: isFifthElementVisible,
				process: openListItem,
			},
			listDialogVariantPrototype{
				id: "back",
				text: trans("back_btn"),
				isActiveFn: isNotTheFirstPage,
				process: moveBack,
			},
			listDialogVariantPrototype{
				id: "fwd",
				text: trans("fwd_btn"),
				isActiveFn: isNotTheLastPage,
				process: moveForward,
			},
		},
	})
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

func getFirstItemId(cahce *listDialogCache) string {
	index := cahce.currentPage * 4
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getSecondItemId(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 1
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getThirdItemId(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 2
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getFourthItemId(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 3
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getFifthItemId(cahce *listDialogCache) string {
	index := cahce.currentPage * 4 + 4
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func addList(additionalId string, data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_enter_list_name"))
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "listname",
	})
	return true
}

func moveForward(additionalId string, data *processing.ProcessData) bool {
	ids, _ := data.Static.Db.GetUserLists(data.UserId)
	itemsCount := len(ids)
	var pagesCount int
	if itemsCount > 2 {
		pagesCount = (itemsCount - 2) / 4 + 1
	} else {
		pagesCount = 1
	}

	currentPage := data.Static.GetUserStateCurrentPage(data.UserId)

	if currentPage + 1 < pagesCount {
		data.Static.SetUserStateCurrentPage(data.UserId, currentPage + 1)
	}
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func moveBack(additionalId string, data *processing.ProcessData) bool {
	currentPage := data.Static.GetUserStateCurrentPage(data.UserId)
	if currentPage > 0 {
		data.Static.SetUserStateCurrentPage(data.UserId, currentPage - 1)
	}
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func openListItem(additionalId string, data *processing.ProcessData) bool {
	id, err := strconv.ParseInt(additionalId, 10, 64)

	if err != nil {
		return false
	}

	if data.Static.Db.IsListBelongsToUser(data.UserId, id) {
		data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("li", id, data.Static))
		return true
	} else {
		return false
	}
}

func (factory *listDialogFactory) createVariants(userId int64, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)
	cache := getListDialogCache(userId, staticData)
	
	row := 1
	col := 0

	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(cache) {
			var text string

			if variant.textFn != nil {
				text = variant.textFn(cache)
			} else {
				text = variant.text
			}

			var additionalId string

			if variant.additionalIdFn != nil {
				additionalId = variant.additionalIdFn(cache)
			}

			variants = append(variants, dialog.Variant{
				Id:   variant.id,
				Text: text,
				AdditionalId: additionalId,
				RowId: row,
			})
			
			col = col + 1
			if col > 1 {
				row = row + 1
				col = 0
			}
		}
	}
	return
}

func getListDialogCache(userId int64, staticData *processing.StaticProccessStructs) (cache *listDialogCache) {

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

func (factory *listDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(additionalId, data)
		}
	}
	return false
}
