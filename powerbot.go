package powerbot

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/thoj/go-ircevent"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (bot *Bot) WriteCode(code int) error {
	fmt.Printf("Sending %v\n", code)
	cmd := exec.Command("digi-rc-switch.py", strconv.Itoa(code))
	stdout, err := cmd.Output()
	fmt.Println(stdout)
	fmt.Println(err)
	return err
}

type Bot struct {
	Name      string
	Channels  []string
	IrcConfig IrcServerConfig
	Con       *irc.Connection
	Commands  map[string][]int
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
	} else if codes, ok := bot.Commands[command+" "+argument]; ok {
		successful_codes := []int{}
		failed_codes := map[int]error{}
		for _, code := range codes {
			if write_err := bot.WriteCode(code); write_err == nil {
				successful_codes = append(successful_codes, code)
			} else {
				failed_codes[code] = write_err
			}
			time.Sleep(time.Second)
		}
		reply := fmt.Sprintf("Sent out codes %v for %v.", successful_codes, command)
		for code, write_err := range failed_codes {
			reply += fmt.Sprintf(" Error writing %v to device: %v.", code, write_err)
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
	IrcServer IrcServerConfig   `yaml:"ircserver"`
	Nick      string            `yaml:"nick"`
	Channels  []string          `yaml:"channels"`
	Commands  map[string][]int  `yaml:"commands"`
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
