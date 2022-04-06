package bot

import (
	"book-club-bot/config"
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	books "google.golang.org/api/books/v1"
)

var (
	BotId        string
	goBot        *discordgo.Session
	booksService *books.Service
)

const (
	baseCommand = "!"
	getVol      = "getBook"
)

func Start() {

	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()
	booksService, err = books.NewService(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)
	goBot.AddHandler(googleBooksHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc)
	<-sc

	goBot.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func googleBooksHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	command, body := parseCommand(m.Content)
	if command == getVol {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Getting book")
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Not a valid command")
	}
	fmt.Println(body)

}

func parseCommand(messageContent string) (parsedCommand string, body string) {
	splitMessage := strings.Split(messageContent, " ")

	match, err := regexp.MatchString("![A-z]+", splitMessage[0])
	if err != nil {
		fmt.Println(err.Error())
		return "", ""
	}

	if match {
		command := splitMessage[0][1:]
		body := strings.Join(splitMessage[1:], " ")
		return command, body
	}

	return
}
