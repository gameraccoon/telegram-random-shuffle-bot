package processing

import (
	"github.com/gameraccoon/telegram-random-shuffle-bot/chat"
	"github.com/gameraccoon/telegram-random-shuffle-bot/database"
	"github.com/gameraccoon/telegram-random-shuffle-bot/dialog"
	"github.com/nicksnyder/go-i18n/i18n"
	"time"
)

type StaticConfiguration struct {
	Language    string
	ExtendedLog bool
}

type AwaitingTextProcessorData struct {
	ProcessorId string
	AdditionalId string
}

type UserState struct {
	awaitingTextProcessor *AwaitingTextProcessorData
	currentPage int
	lastMessages []int64
}

type StaticProccessStructs struct {
	Chat       chat.Chat
	Db         *database.Database
	Timers     map[int64]time.Time
	Config     *StaticConfiguration
	Trans      i18n.TranslateFunc
	MakeDialogFn func(string, int64, *StaticProccessStructs)*dialog.Dialog
	userStates map[int64]UserState
}

func (staticData *StaticProccessStructs) Init() {
	staticData.userStates = make(map[int64]UserState)
}

func (staticData *StaticProccessStructs) SetUserStateTextProcessor(userId int64, proessor *AwaitingTextProcessorData) {
	state := staticData.userStates[userId]
	state.awaitingTextProcessor = proessor
	staticData.userStates[userId] = state
}

func (staticData *StaticProccessStructs) GetUserStateTextProcessor(userId int64) *AwaitingTextProcessorData {
	if state, ok := staticData.userStates[userId]; ok {
		return state.awaitingTextProcessor
	} else {
		return nil
	}
}

func (staticData *StaticProccessStructs) SetUserStateCurrentPage(userId int64, page int) {
	state := staticData.userStates[userId]
	state.currentPage = page
	staticData.userStates[userId] = state
}

func (staticData *StaticProccessStructs) GetUserStateCurrentPage(userId int64) int {
	if state, ok := staticData.userStates[userId]; ok {
		return state.currentPage
	} else {
		return 0
	}
}
