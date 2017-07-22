package main

import (
	//	"bytes"
	//	"fmt"
	"github.com/gameraccoon/telegram-random-shuffle-bot/database"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialogFactories"
	"github.com/gameraccoon/telegram-random-shuffle-bot/processing"
	//"github.com/gameraccoon/telegram-random-shuffle-bot/telegramChat"
	//"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"strconv"
	"strings"
)

type ProcessorFunc func(*processing.ProcessData, *dialogFactories.DialogManager)

type ProcessorFuncMap map[string]ProcessorFunc

func sendEditingGuide(data *processing.ProcessData, dialogManager *dialogFactories.DialogManager) {
	dialog := dialogManager.MakeDialog("ed", data)
	if dialog != nil {
		data.Static.Chat.SendDialog(dialog, data.ChatId)
	}
}

func startCommand(data *processing.ProcessData, dialogManager *dialogFactories.DialogManager) {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("hello_message"))
}

func setVariants(db *database.Database, id int64, text *string) (ok bool) {
	//variants := strings.Split(*text, "\n")
	//db.SetQuestionVariants(questionId, variants)
	return true
}

func makeUserCommandProcessors() ProcessorFuncMap {
	return map[string]ProcessorFunc{
		"start": startCommand,
	}
}

func processCommandByProcessors(data *processing.ProcessData, processors *ProcessorFuncMap, dialogManager *dialogFactories.DialogManager) bool {
	processor, ok := (*processors)[data.Command]
	if ok {
		processor(data, dialogManager)
	}

	return ok
}

func processCommand(data *processing.ProcessData, dialogManager *dialogFactories.DialogManager, processors *ProcessorFuncMap) {
	// process dialogs
	ids := strings.Split(data.Command, "_")
	if len(ids) >= 2 {
		processed := dialogManager.ProcessVariant(ids[0], ids[1], data)
		if processed {
			return
		}
	}

	// process static command
	processed := processCommandByProcessors(data, processors, dialogManager)
	if processed {
		return
	}

	// if we here it means that no command was processed
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("warn_unknown_command"))
}

func processPlainMessage(data *processing.ProcessData, dialogManager *dialogFactories.DialogManager) {
	if _, ok := data.Static.UserStates[data.UserId]; ok {
		success := dialogManager.ProcessText(data)

		if !success {
			data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("warn_unknown_command"))
		}
	} else {
		data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("warn_unknown_command"))
	}
}

func processUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogFactories.DialogManager, processors *ProcessorFuncMap) {
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
