package main

import (
	//	"bytes"
	//	"fmt"
	// "github.com/gameraccoon/telegram-random-shuffle-bot/database"
	// "github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactories"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogManager"
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	//"github.com/gameraccoon/telegram-random-shuffle-bot/telegramChat"
	//"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"strconv"
	"strings"
)

type ProcessorFunc func(*processing.ProcessData)

type ProcessorFuncMap map[string]ProcessorFunc

func startCommand(data *processing.ProcessData) {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("hello_message"))
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
}

func makeUserCommandProcessors() ProcessorFuncMap {
	return map[string]ProcessorFunc{
		"start": startCommand,
	}
}

func processCommandByProcessors(data *processing.ProcessData, processors *ProcessorFuncMap) bool {
	processor, ok := (*processors)[data.Command]
	if ok {
		processor(data)
	}

	return ok
}

func processCommand(data *processing.ProcessData, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	// drop any text processors for the case wi will process a command
	data.Static.SetUserStateTextProcessor(data.UserId, nil)
	// process dialogs
	ids := strings.Split(data.Command, "_")
	if len(ids) >= 2 {
		dialogId := ids[0]
		variantId := ids[1]
		var additionalId string
		if len(ids) > 2 {
			additionalId = ids[2]
		}

		processed := dialogManager.ProcessVariant(dialogId, variantId, additionalId, data)
		if processed {
			return
		}
	}

	// process static command
	processed := processCommandByProcessors(data, processors)
	if processed {
		return
	}

	// if we here that means that no command was processed
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("warn_unknown_command"))
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Static))
}

func processPlainMessage(data *processing.ProcessData, dialogManager *dialogManager.DialogManager) {
	success := dialogManager.ProcessText(data)

	if !success {
		data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("warn_unknown_command"))
	}
}

func processUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	data := processing.ProcessData{
		Static: staticData,
		ChatId: update.Message.Chat.ID,
		UserId: staticData.Db.GetUserId(update.Message.Chat.ID),
	}

	message := update.Message.Text

	if strings.HasPrefix(message, "/") {
		commandLen := strings.Index(message, " ")
		if commandLen != -1 {
			data.Command = message[1:commandLen]
			data.Message = message[commandLen+1:]
		} else {
			data.Command = message[1:]
		}

		processCommand(&data, dialogManager, processors)
	} else {
		data.Message = message
		processPlainMessage(&data, dialogManager)
	}
}
