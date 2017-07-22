package dialogFactories

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
)

type DialogManager struct {
	dialogs map[string]*DialogFactory
	textProcessors textInputProcessorManager
}

func (dialogManager *DialogManager) RegisterDialogFactory(id string, dialogFactory *DialogFactory) {
	if dialogManager.dialogs == nil {
		dialogManager.dialogs = make(map[string]*DialogFactory)
	}

	dialogManager.dialogs[id] = dialogFactory
	dialogFactory.id = id
}

func (dialogManager *DialogManager) InitTextProcessors() {
	dialogManager.textProcessors = getTextInputProcessorManager()
}

func (dialogManager *DialogManager) MakeDialog(dialogId string, data *processing.ProcessData) (dialog *dialog.Dialog) {
	factory := dialogManager.getDialogFactory(dialogId)
	if factory != nil {
		dialog = factory.MakeDialog(data)
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

func (dialogManager *DialogManager) getDialogFactory(id string) *DialogFactory {
	dialogFactory, ok := dialogManager.dialogs[id]
	if ok && dialogFactory != nil {
		return dialogFactory
	} else {
		return nil
	}
}
