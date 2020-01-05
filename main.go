package main

import (
	"github.com/romanornr/Bitmex-referral-analyzer/bots"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"github.com/spf13/viper"
)

func init() {
	config.GetViperConfig()
}

func main() {
	telegramBot := bots.NewTelegramBot(viper.GetString("telegram.token"))
	telegramBot.Update()
}