package main

import (
	"fmt"
	"github.com/solarkennedy/powerbot"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	fmt.Printf("DEBUG: Config:: %+v", config)
	bot := powerbot.Bot{
		Name:      config.Nick,
		Channels:  config.Channels,
		IrcConfig: config.IrcServer,
		Commands:  config.Commands}
	arg := os.Args[1]
	fmt.Printf("Executing CLI commadnd out of %+v", arg)
	code, err := strconv.Atoi(arg)
	if err == nil {
		write_err := bot.WriteCode(code)
		if write_err == nil {
			fmt.Printf("Sent out code %v", code)
			os.Exit(0)
		} else {
			fmt.Printf("Error sending out code %v:", code)
			fmt.Printf("%v", write_err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error parsing the code to send:")
		fmt.Printf("%v", err)
		os.Exit(3)
	}
}
