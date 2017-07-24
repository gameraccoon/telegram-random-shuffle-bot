package dialogManager

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactory"
)

type DialogManager struct {
	dialogs map[string]dialogFactory.DialogFactory
	textProcessors textInputProcessorManager
}

func (dialogManager *DialogManager) RegisterDialogFactory(id string, factory dialogFactory.DialogFactory) {
	if dialogManager.dialogs == nil {
		dialogManager.dialogs = make(map[string]dialogFactory.DialogFactory)
	}

	dialogManager.dialogs[id] = factory
}

func (dialogManager *DialogManager) InitTextProcessors() {
	dialogManager.textProcessors = getTextInputProcessorManager()
}

func (dialogManager *DialogManager) MakeDialog(dialogId string, data *processing.ProcessData) (dialog *dialog.Dialog) {
	factory := dialogManager.getDialogFactory(dialogId)
	if factory != nil {
		dialog = factory.MakeDialog(data)
		dialog.Id = dialogId
	}
	return
}

func (dialogManager *DialogManager) ProcessVariant(dialogId string, variantId string, data *processing.ProcessData) (processed bool) {
	factory := dialogManager.getDialogFactory(dialogId)
	if factory != nil {
		processed = factory.ProcessVariant(variantId, data)
	}
	return
}

func (dialogManager *DialogManager) ProcessText(data *processing.ProcessData) bool {
	return dialogManager.textProcessors.processText(data)
}

func (dialogManager *DialogManager) getDialogFactory(id string) dialogFactory.DialogFactory {
	dialogFactory, ok := dialogManager.dialogs[id]
	if ok && dialogFactory != nil {
		return dialogFactory
	} else {
		return nil
	}
}
