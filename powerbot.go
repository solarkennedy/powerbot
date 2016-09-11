package main

import (
	"encoding/binary"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	"os"
	"strings"
)

func codeToBytes(code int) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(code))
	fmt.Printf("DEBUG: code: %v, bytes: %v\n", code, bs)
	return bs
}

func writeCode(code int, serialport string) {
	bs := codeToBytes(code)
	fmt.Printf("Sending %v\n", bs)
	fi, err := os.OpenFile("/dev/ttyACM0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic(err)
	}
	_, err2 := fi.Write(bs)
	if err2 != nil {
		panic(err2)
	}
	fi.Sync()
	fi.Close()
}

type Bot struct {
	Name       string
	Room       string
	Server     string
	Port       int
	Con        *irc.Connection
	SerialPort string
}

func (bot *Bot) Address() string {
	return fmt.Sprintf("%s:%d", bot.Server, bot.Port)
}

func (bot *Bot) ParseAndReply(channel string, msg string) {
	bot.Con.Privmsg(channel, msg)
}

func (bot *Bot) Run() {
	bot.Con = irc.IRC(bot.Name, bot.Name)
	err := bot.Con.Connect(bot.Address())
	if err != nil {
		log.Fatal("Couldn't connect to %s: %s", bot.Address, err)
	}
	bot.Con.AddCallback("001", func(e *irc.Event) {
		bot.Con.Join(bot.Room)
	})
	bot.Con.AddCallback("PRIVMSG", func(e *irc.Event) {
		channel := e.Arguments[0]
		msg := e.Arguments[1]
		if strings.HasPrefix(msg, bot.Name) {
			bot.ParseAndReply(channel, msg)
		}
	})
	bot.Con.Loop()
}

func main() {
	bot := Bot{
		Name: "powerbot", Room: "#test", Server: "archive.local", Port: 6667,
		SerialPort: "/dev/ttyACM0"}
	bot.Run()
}
