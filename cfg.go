package main

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	BotToken        string `json:"botToken"`
	ChannelUsername string `json:"channelUsername"`
	DataDIR         string `json:"dataDir"`
}

func loadConfig(confPath string) *config {
	conf := &config{}
	file, err := os.Open(confPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	doc := json.NewDecoder(file)
	doc.Decode(&conf)

	return conf
}
