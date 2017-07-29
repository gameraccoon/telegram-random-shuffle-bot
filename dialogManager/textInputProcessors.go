package dialogManager

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"strconv"
	"strings"
)

type textProcessorFunc func(string, *processing.ProcessData) bool

type textProcessorsMap map[string]textProcessorFunc

type textInputProcessorManager struct {
	processors textProcessorsMap
}

func getTextInputProcessorManager() textInputProcessorManager {
	return textInputProcessorManager {
		processors : textProcessorsMap {
			"listname" : processAddList,
			"addlistitems" : processAddItemsToList,
		},
	}
}

func (textProcessorManager *textInputProcessorManager) processText(data *processing.ProcessData) bool {
	textProcessor := data.Static.GetUserStateTextProcessor(data.UserId)
	if textProcessor != nil {
		processor, ok := textProcessorManager.processors[textProcessor.ProcessorId]
		if ok {
			return processor(textProcessor.AdditionalId, data)
		}
	}
	return false
}

func processAddList(additionalId string, data *processing.ProcessData) bool {
	newListId := data.Static.Db.CreateList(data.UserId, data.Message)
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "addlistitems",
		AdditionalId: strconv.FormatInt(newListId, 10),
	})
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_wait_items"))
	return true
}

func processAddItemsToList(additionalId string, data *processing.ProcessData) bool {
	items := strings.Split(data.Message, "\n")
	
	listId, err := strconv.ParseInt(additionalId, 10, 64)
	if err != nil {
		return false
	}
	
	data.Static.Db.AddItemsToList(listId, items)
	
	data.Static.SetUserStateTextProcessor(data.UserId, nil)
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_items_added"))
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("li", listId, data.Static))
	return true
}
