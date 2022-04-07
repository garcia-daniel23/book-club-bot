package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token             string
	BotPrefix         string
	GoogleBooksAPIKey string

	config *configStruct
)

type configStruct struct {
	Token             string `json : "Token"`
	BotPrefix         string `json : "BotPrefix"`
	GoogleBooksAPIKey string `json : GoogleBooksAPIKey`
}

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Successfully read config file.")

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	GoogleBooksAPIKey = config.GoogleBooksAPIKey

	return nil
}
