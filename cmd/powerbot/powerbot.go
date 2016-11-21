package main

import (
	"fmt"
	"github.com/solarkennedy/powerbot"
	"io/ioutil"
	"log"
)

func main() {
	filename := "powerbot.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var config powerbot.Config
	err = config.Parse(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%+v", config)
	bot := powerbot.Bot{
		Name:       config.Nick,
		Channels:   config.Channels,
		IrcConfig:  config.IrcServer,
		SerialPort: config.SerialPort,
		Commands:   config.Commands}
	bot.Run()
}
