package bot

import (
	"book-club-bot/config"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	books "google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

const (
	getBookCmd = "!getbook"
)

var (
	BotId        string
	goBot        *discordgo.Session
	booksService *books.Service

	commands map[string]func(messageBody string) (response string)
)

//Return list of books matching users input from GoogleBooks API
func searchForBooks(searchTitle string) ([]*books.Volume, error) {
	fmt.Println(searchTitle)

	volumesListCall := booksService.Volumes.List(searchTitle)
	volumes, err := volumesListCall.Do()
	if err != nil {
		return nil, err
	}

	return volumes.Items, nil
}

//Use fuzzy wuzzy search on returned list to find the best matching object returned.
func findBestBook(searchTitle string) (book *books.Volume, err error) {
	var bestMatchedBook *books.Volume
	bestMatchValue := 0
	books, err := searchForBooks(searchTitle)
	if err != nil {
		return nil, errors.New("error in searchForBooks: " + err.Error())
	}

	for _, book := range books {
		bookTitle := book.VolumeInfo.Title + " " + book.VolumeInfo.Subtitle

		fmt.Println("Title: " + bookTitle)
		fmt.Printf("Simple Ratio option: %v\n", fuzzy.Ratio(searchTitle, bookTitle))
		fmt.Printf("Partial Ratio option: %v\n", fuzzy.PartialRatio(searchTitle, bookTitle))
		fmt.Printf("TokenSet option: %v\n", fuzzy.TokenSetRatio(searchTitle, bookTitle))
		fmt.Printf("TokenSort option: %v\n", fuzzy.TokenSortRatio(searchTitle, bookTitle))
		matchValue := fuzzy.Ratio(searchTitle, bookTitle)
		if matchValue > bestMatchValue {
			bestMatchedBook = book
			bestMatchValue = matchValue
		}

	}

	if bestMatchedBook != nil {
		return bestMatchedBook, nil
	}

	return bestMatchedBook, errors.New("no book was found")
}

func getBookCommand(searchTitle string) (response string) {

	fmt.Println("executing book command")

	book, err := findBestBook(searchTitle)
	if err != nil {
		fmt.Println(errors.New("getBookCommand failed: " + err.Error()))
		return
	}

	previewLink := book.VolumeInfo.PreviewLink
	title := book.VolumeInfo.Title
	if book.VolumeInfo.Subtitle != "" {
		title += " " + book.VolumeInfo.Subtitle
	}

	author := book.VolumeInfo.Authors
	description := book.VolumeInfo.Description
	// thumbnail := book.VolumeInfo.ImageLinks.Thumbnail

	res := fmt.Sprintf("%v"+
		"\nTitle: %v"+
		"\nAuthor: %v"+
		"\nDescription: %v", previewLink, title, strings.Join(author, " "), description)

	return res
}

func Start() {
	//Initialize Commands
	commands = map[string]func(messageBody string) (response string){
		getBookCmd: getBookCommand,
	}

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	booksService, err = books.NewService(ctx, option.WithAPIKey(config.GoogleBooksAPIKey))
	if err != nil {
		fmt.Println(err.Error())
	}

	BotId = u.ID

	goBot.AddHandler(commandHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running! Press CTRL-C to exit.")

	<-ctx.Done()
	stop()

	fmt.Println("Shutting down bot")

	goBot.Close()
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	splitMessage := strings.Split(m.Content, " ")

	match, err := regexp.MatchString("![A-z]+", splitMessage[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !match {
		return
	}

	commandName := splitMessage[0] //Grab the command name ex: '!getBook <command body>' where !getBook is the command
	commandBody := strings.Join(splitMessage[1:], " ")

	if val, ok := commands[commandName]; ok {
		response := val(commandBody)
		_, _ = s.ChannelMessageSend(m.ChannelID, response)
	}
}
