package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogManager"
	"strconv"
	"strings"
)

func GetTextInputProcessorManager() dialogManager.TextInputProcessorManager {
	return dialogManager.TextInputProcessorManager {
		Processors : dialogManager.TextProcessorsMap {
			"listname" : processAddList,
			"addlistitems" : processAddItemsToList,
		},
	}
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
