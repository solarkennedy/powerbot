package main

import (
	"encoding/binary"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func codeToBytes(code int) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(code))
	fmt.Printf("DEBUG: code: %v, bytes: %v\n", code, bs)
	return bs
}

func (bot *Bot) WriteCode(code int) error {
	bs := codeToBytes(code)
	fmt.Printf("Sending %v to %s\n", bs, bot.SerialPort)
	fi, err := os.OpenFile(bot.SerialPort, os.O_WRONLY, os.ModeDevice)
	if err != nil {
		return err
	}
	_, err2 := fi.Write(bs)
	if err2 != nil {
		return err2
	}
	return nil
}

type Bot struct {
	Name       string
	Channels   []string
	IrcConfig  IrcServerConfig
	Con        *irc.Connection
	SerialPort string
	Commands   map[string]int
}

func (bot *Bot) Address() string {
	return fmt.Sprintf("%s:%d", bot.IrcConfig.Hostname, bot.IrcConfig.Port)
}

func ExtractCommandAndArgument(msg string, name string) (command string, argument string) {
	regex := fmt.Sprintf(`^%s[:]?.* (\w+) (\w+)`, name)
	command_regexp := regexp.MustCompile(regex)
	matches := command_regexp.FindStringSubmatch(msg)
	if len(matches) != 3 {
		command = "unknown"
		argument = "unknown"
	} else {
		command = matches[1]
		argument = matches[2]
	}
	return
}

func (bot *Bot) ParseAndReply(channel string, msg string, user string) {
	command, argument := ExtractCommandAndArgument(msg, bot.Name)
	if command == "code" {
		code, err := strconv.Atoi(argument)
		if err == nil {
			write_err := bot.WriteCode(code)
			var reply string
			if write_err == nil {
				reply = fmt.Sprintf("Sent out code %v", code)
			} else {
				reply = fmt.Sprintf("Error writing to device: %v", write_err)
			}
			bot.Con.Privmsg(channel, reply)
			return
		} else {
			reply := fmt.Sprintf("%v doesn't look like a valid code", argument)
			bot.Con.Privmsg(channel, reply)
			return
		}
	} else if code, ok := bot.Commands[command+" "+argument]; ok {
		write_err := bot.WriteCode(code)
		var reply string
		if write_err == nil {
			reply = fmt.Sprintf("Sent out code %v for %v", code, command)
		} else {
			reply = fmt.Sprintf("Error writing to device: %v", write_err)
		}
		bot.Con.Privmsg(channel, reply)
		return
	} else {
		bot.Con.Privmsg(user, fmt.Sprintf("'%v %v' is not a valid command", command, argument))
		bot.Con.Privmsg(user, "To send raw codes:")
		bot.Con.Privmsg(user, "    powerbot code 1234")
		bot.Con.Privmsg(user, "Or one of the configured commands:")
		all_commands := fmt.Sprintf("    %v", bot.ListCommands())
		bot.Con.Privmsg(user, all_commands)
		return
	}
}

func (bot *Bot) ListCommands() []string {
	keys := make([]string, 0, len(bot.Commands))
	for k, _ := range bot.Commands {
		keys = append(keys, k)
	}
	return keys
}

func (bot *Bot) Run() {
	bot.Con = irc.IRC(bot.Name, bot.Name)
	bot.Con.UseTLS = bot.IrcConfig.SSL
	bot.Con.Password = bot.IrcConfig.Password
	bot.Con.VerboseCallbackHandler = true
	bot.Con.Debug = true
	err := bot.Con.Connect(bot.Address())
	if err != nil {
		log.Fatal("Couldn't connect to %s: %s", bot.Address(), err)
	}
	bot.Con.Nick(bot.Name)
	log.Printf("Connected to %v as %v", bot.Address(), bot.Name)
	for _, channel := range bot.Channels {
		log.Printf("Joining %v", channel)
		bot.Con.AddCallback("001", func(e *irc.Event) {
			bot.Con.Join(channel)
		})
	}
	bot.Con.AddCallback("PRIVMSG", func(e *irc.Event) {
		channel := e.Arguments[0]
		msg := e.Arguments[1]
		if strings.HasPrefix(msg, bot.Name) {
			bot.ParseAndReply(channel, msg, e.Nick)
		}
	})
	bot.Con.Loop()
}

type IrcServerConfig struct {
	Hostname string `yaml:"hostname"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	SSL      bool   `yaml:"ssl"`
}

func (c *IrcServerConfig) UnmarshalYAML(b []byte) error {
	return yaml.Unmarshal(b, c)
}

type Config struct {
	SerialPort string          `yaml:"serialport"`
	IrcServer  IrcServerConfig `yaml:"ircserver"`
	Nick       string          `yaml:"nick"`
	Channels   []string        `yaml:"channels"`
	Commands   map[string]int  `yaml:"commands"`
}

func (c *Config) Parse(data []byte) error {
	err := yaml.Unmarshal([]byte(data), c)
	if err != nil {
		return err
	}
	if c.IrcServer.Hostname == "" {
		log.Fatalf("error: ircserver hostname not set")
	}
	return nil
}

func main() {
	filename := "powerbot.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var config Config
	err = config.Parse(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%+v", config)
	bot := Bot{
		Name:       config.Nick,
		Channels:   config.Channels,
		IrcConfig:  config.IrcServer,
		SerialPort: config.SerialPort,
		Commands:   config.Commands}
	bot.Run()
}
