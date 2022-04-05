package main

import (
	"book-club-bot/bot"
	"book-club-bot/config"
	"fmt"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	return
}
