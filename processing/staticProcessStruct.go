package processing

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/chat"
	"github.com/gameraccoon/telegram-random-shuffle-bot/database"
	"github.com/nicksnyder/go-i18n/i18n"
	"time"
)

type StaticConfiguration struct {
	Language    string
	ExtendedLog bool
}

type UserState struct {
	AwaitingTextProcessor *string
}

type StaticProccessStructs struct {
	Chat       chat.Chat
	Db         *database.Database
	UserStates map[int64]UserState
	Timers     map[int64]time.Time
	Config     *StaticConfiguration
	Trans      i18n.TranslateFunc
}

func (staticData *StaticProccessStructs) SetUserStateTextProcessor(userId int64, proessor string) {
	state := staticData.UserStates[userId]
	state.AwaitingTextProcessor = &proessor
	staticData.UserStates[userId] = state
}

func (staticData *StaticProccessStructs) ClearUserStateTextProcessor(userId int64) {
	state := staticData.UserStates[userId]
	state.AwaitingTextProcessor = nil
	staticData.UserStates[userId] = state
}

func (staticData *StaticProccessStructs) GetUserStateTextProcessor(userId int64) *string {
	if state, ok := staticData.UserStates[userId]; ok {
		return state.AwaitingTextProcessor
	} else {
		return nil
	}
}
