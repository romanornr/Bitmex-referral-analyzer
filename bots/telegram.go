package bots

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/romanornr/Bitmex-referral-analyzer/bmex"
	"log"
	"strings"
	"time"
)
const PARSEMODE = "html"

type TelegramBot struct {
	api string
	bot *tgbotapi.BotAPI
}

func NewTelegramBot(api string) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		log.Panicf("error new telegram bot: %s\n", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &TelegramBot{api:api, bot:bot}
}

func (telegram TelegramBot) Update() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := telegram.bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("failed to receive message from channel: %s\n")
	}

	for update := range updates {
		if update.Message == nil { // ignore any non messages
			continue
		}

		s, err := NewSession(&update, telegram.bot)
		if err != nil {
			log.Printf("failed creating session: %s\n", s)
		}
		go s.Handle()
	}
}

type Session struct {
	update  *tgbotapi.Update
	command *Command
	bot     *tgbotapi.BotAPI
}

type Command struct {
	Command string
	Params  []string
}

func NewSession(u *tgbotapi.Update, b *tgbotapi.BotAPI) (s *Session, err error) {
	s = &Session{update: u, bot: b}

	if !s.checkMessageTimeVerify() {
		return s, fmt.Errorf("check message time verify failed")
	}

	return s, nil
}

// check the time difference between the message and the current time.
// if the time difference is too big, return false (don't sent anything)
// this is for if the bot goes down and hundreds of commands are still in queue.
func (s Session) checkMessageTimeVerify() bool {
	messageTime := int64(s.update.Message.Date)
	now := time.Now().Unix()
	difference := now - messageTime
	if difference > 60 {
		return false
	}
	return true
}

// Handle runs known commands for the user session
func (s *Session) Handle() {
	session := CommandParser(s)
	botUsername := "@" + session.bot.Self.UserName

	m := map[string]interface{}{
		"/week" + botUsername:   getWeeklyEarnings,
		"/months" + botUsername: getMonthlyEarnings,
	}

	if _, ok := m[s.command.Command]; ok {
		m[s.command.Command].(func(*Session))(session)
	}
}

func CommandParser(s *Session) *Session {
	var command string

	//add username to the command if message does not contain username
	//example /index via into /@btsebot
	botUsername := s.bot.Self.UserName
	command = strings.Split(s.update.Message.Text, " ")[0]
	if !strings.Contains(s.update.Message.Text, botUsername) {
		command = strings.Split(s.update.Message.Text, " ")[0] + "@" + botUsername
	}
	temp := strings.TrimLeft(s.update.Message.Text, botUsername)
	params := strings.Split(temp, " ")[1:]
	s.command = &Command{Command: command, Params: params}
	return s
}

func getMonthlyEarnings(s *Session) {
	err, tx := bmex.LoadWalletHistory()
	if err != nil {
		log.Printf("Failed to load wallet transactions: %s\n", err)
		message := fmt.Sprintf("Failed to load wallet transactions: %s\n", err)
		sendMessage(s, message)
	}
	bmex.ReferralEarning(tx)

	//message := "<code>"
	//stats := bmex.ReferralEarning(tx)
	//for i := len(stats.Stat) - 1; i >= 0; i-- {
	//	message += fmt.Sprintf("%s\t\t\t%s\n", stats.Stat[i].Date, stats.Stat[i].Dollar)
	//}
	//message += fmt.Sprintf("\nTotal BTC: %s \t Total Dollar: %s</code>",stats.TotalBtc, stats.TotalDollar)
	//
	//sendMessage(s, message)
}

func getWeeklyEarnings(s *Session) {
	err, tx := bmex.LoadWalletHistory()
	if err != nil {
		log.Printf("Failed to load wallet transactions: %s\n", err)
		message := fmt.Sprintf("Failed to load wallet transactions: %s\n", err)
		sendMessage(s, message)
	}
	message := "<code>"
	stats := bmex.WeeklyEarnings(tx)

	for i := len(stats.Stat) - 1; i >= 0; i-- {
		message += fmt.Sprintf("%s\t\t\t%s\n", stats.Stat[i].Date, stats.Stat[i].Dollar)
	}
	message += fmt.Sprintf("\nTotal BTC: %s \t Total Dollar: %s</code>",stats.TotalBtc, stats.TotalDollar)

	sendMessage(s, message)
}

func sendMessage(s *Session, message string) {
	msg := tgbotapi.NewMessage(s.update.Message.Chat.ID, message)
	msg.ReplyToMessageID = s.update.Message.MessageID
	msg.ParseMode = PARSEMODE
	s.bot.Send(msg)
}
