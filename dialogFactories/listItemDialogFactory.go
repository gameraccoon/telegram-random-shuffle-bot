package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
	"fmt"
	"math/rand"
	"strings"
	"strconv"
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
				id: "delrand",
				text: trans("delete_and_reroll"),
				process: deleteAndGetRandom,
			},
			listItemVariantPrototype{
				id: "rand",
				text: trans("reroll"),
				process: getRandom,
			},
			listItemVariantPrototype{
				id: "mix",
				text: trans("get_shuffled"),
				process: getShuffled,
			},
			listItemVariantPrototype{
				id: "add",
				text: trans("add_btn"),
				process: addItems,
			},
			listItemVariantPrototype{
				id: "del",
				text: trans("remove_list"),
				process: removeList,
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
	listId, err := strconv.ParseInt(additionalId, 10, 64)
	
	if err != nil {
		return false
	}
	
	ids, _ := data.Static.Db.GetListItems(listId)
	
	if len(ids) > 0 {
		choosenId := ids[rand.Int63n(int64(len(ids)))]
		data.Static.Db.SetLastItem(listId, choosenId)
	}
	
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("li", listId, data.Static))
	return true
}

func deleteAndGetRandom(additionalId string, data *processing.ProcessData) bool {
	listId, err := strconv.ParseInt(additionalId, 10, 64)
	
	if err != nil {
		return false
	}

	lastItem, _ := data.Static.Db.GetLastItem(listId)
	if lastItem != -1 {
		data.Static.Db.RemoveItem(lastItem)
		data.Static.Db.SetLastItem(listId, -1)
	}
	
	ids, _ := data.Static.Db.GetListItems(listId)
	
	if len(ids) > 0 {
		choosenId := ids[rand.Int63n(int64(len(ids)))]
		data.Static.Db.SetLastItem(listId, choosenId)
	} else {
		data.Static.Db.SetLastItem(listId, -1)
	}
	
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("li", listId, data.Static))
	return true
}

func getShuffled(additionalId string, data *processing.ProcessData) bool {
	listId, err := strconv.ParseInt(additionalId, 10, 64)
	
	if err != nil {
		return false
	}
	
	_, texts := data.Static.Db.GetListItems(listId)
	
	for i := range texts {
    j := rand.Intn(i + 1)
    texts[i], texts[j] = texts[j], texts[i]
	}
	
	data.Static.Chat.SendMessage(data.ChatId, strings.Join(texts[:], "\n"))
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("li", listId, data.Static))
	return true
}

func addItems(additionalId string, data *processing.ProcessData) bool {
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "addlistitems",
		AdditionalId: additionalId,
	})
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_list_added"))
	return true
}

func removeList(additionalId string, data *processing.ProcessData) bool {
	listId, err := strconv.ParseInt(additionalId, 10, 64)
	
	if err != nil {
		return false
	}
	
	data.Static.Db.DeleteList(listId)
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func backToLists(additionalId string, data *processing.ProcessData) bool {
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
	return true
}

func (factory *listItemDialogFactory) getListItemDialogText(listId int64, staticData *processing.StaticProccessStructs) string {
	id, text := staticData.Db.GetLastItem(listId)
	ids, _ := staticData.Db.GetListItems(listId)
	countText := strconv.FormatInt(int64(len(ids)), 10)
	if id != -1 {
		return fmt.Sprintf("%s\n%s\n%s", staticData.Db.GetListName(listId), countText, text)
	} else {
		return fmt.Sprintf("%s\n%s", staticData.Db.GetListName(listId), countText)
	}
}

func (factory *listItemDialogFactory) createVariants(listId int64, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)
	
	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:   variant.id,
			Text: variant.text,
			AdditionalId: strconv.FormatInt(listId, 10),
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
