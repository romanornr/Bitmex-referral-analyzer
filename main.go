package main

import (
	"fmt"
	"github.com/romanornr/Bitmex-referral-analyzer/bots"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"github.com/spf13/viper"
)

func init() {
	config.GetViperConfig()
}

func main() {
	fmt.Println(viper.GetString("telegram.token"))
	telegramBot := bots.NewTelegramBot(viper.GetString("telegram.token"))
	telegramBot.Update()
}