package main

import (
	"encoding/binary"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	"os"
	"strconv"
	"strings"
)

func codeToBytes(code int) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(code))
	fmt.Printf("DEBUG: code: %v, bytes: %v\n", code, bs)
	return bs
}

func (bot *Bot) WriteCode(code int) {
	bs := codeToBytes(code)
	fmt.Printf("Sending %v\n", bs)
	fi, err := os.OpenFile(bot.SerialPort, os.O_WRONLY, os.ModeDevice)
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

func ExtractCommandAndArgument(msg string) (command string, argument string) {
	command = "code"
	argument = "95500"
	return
}

func (bot *Bot) ParseAndReply(channel string, msg string, user string) {
	command, argument := ExtractCommandAndArgument(msg)
	if command == "code" {
		code, err := strconv.Atoi(argument)
		if err == nil {
			reply := fmt.Sprintf("Sent out code %v", code)
			bot.WriteCode(code)
			bot.Con.Privmsg(channel, reply)
			return
		} else {
			reply := fmt.Sprintf("%v doesn't look like a valid code", argument)
			bot.Con.Privmsg(channel, reply)
			return
		}
	} else {
		bot.Con.Privmsg(user, fmt.Sprintf("%v is not a valid command", command))
		bot.Con.Privmsg(user, "Try something like:")
		bot.Con.Privmsg(user, "powerbot code 1234")
		return
	}
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
			bot.ParseAndReply(channel, msg, e.Nick)
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
