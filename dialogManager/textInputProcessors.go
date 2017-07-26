package dialogManager

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
)

type textProcessorFunc func(*processing.ProcessData) bool

type textProcessorsMap map[string]textProcessorFunc

type textInputProcessorManager struct {
	processors textProcessorsMap
}

func getTextInputProcessorManager() textInputProcessorManager {
	return textInputProcessorManager {
		processors : textProcessorsMap {
			"list_name" : processSetVariantsContent,
		},
	}
}

func (textProcessorManager *textInputProcessorManager) processText(data *processing.ProcessData) bool {
	textProcessor := data.Static.GetUserStateTextProcessor(data.UserId)
	if textProcessor != nil {
		processor, ok := textProcessorManager.processors[*textProcessor]
		if ok {
			return processor(data)
		}
	}
	return false
}

func processSetVariantsContent(data *processing.ProcessData) bool {
	data.Static.Db.CreateList(data.UserId, data.Message)
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_variants_is_set"))
	return true
}
