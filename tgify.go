package main

import (
	"bufio"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type tgConf struct {
	Token   string  `yaml:"token"`
	ChatIds []int64 `yaml:"chatIds"`
}

func readTelegramConfig() (tgConf, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Get homer dir error:", err)
		return tgConf{}, err
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(home, ".tgify/config.yaml"))
	if err != nil {
		fmt.Println("Read tgify config error:", err)
		return tgConf{}, err
	}

	var cfg tgConf
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		fmt.Println("Unmarshal error:", err)
		return tgConf{}, err
	}

	return cfg, nil
}

type tgIface struct {
	Bot *tgbotapi.BotAPI
	Cfg *tgConf
}

func (tgi tgIface) Printf(format string, args ...interface{}) {
	textMsg := fmt.Sprintf(format, args...)

	for _, chatId := range tgi.Cfg.ChatIds {
		msg := tgbotapi.NewMessage(chatId, textMsg)
		if _, err := tgi.Bot.Send(msg); err != nil {
			fmt.Println("Send msg failed:", err)
		}
	}
	return
}

func (tgi tgIface) SignalHandler() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tgi.Bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := strings.ToLower(update.Message.Text)
		if msg == "s" || msg == "stop" {
			tgi.Bot.StopReceivingUpdates()
			os.Exit(1)
		}
	}
}

func readLines(src *bufio.Reader, prefix string, dstIface *tgIface) {
	for {
		line, _, err := src.ReadLine()
		if err == io.EOF {
			if err != io.EOF {
				fmt.Println("Readline error: %s", err)
			}
			break
		}
		dstIface.Printf("%s%s \n", prefix, line)
	}
}

func main() {

	cfg, err := readTelegramConfig()
	if err != nil {
		fmt.Println("Please create valid tgify config.")
		return
	}

	streamRun := len(os.Args) == 1
	if !streamRun && (len(os.Args) < 3 && os.Args[1] != "--args") {
		fmt.Println("Invalid args:\n\ttgify --args executable-file [inferior-arguments ...]")
		return
	}

	tgi := &tgIface{
		Bot: func() *tgbotapi.BotAPI {
			bot, err := tgbotapi.NewBotAPI(cfg.Token)
			if err != nil {
				fmt.Println("tgbotapi.NewBotAPI error:", err)
				os.Exit(-1)
				return nil
			}
			return bot
		}(),
		Cfg: &cfg,
	}

	go tgi.SignalHandler()

	if streamRun {
		readLines(bufio.NewReader(os.Stdin), "", tgi)
		return
	}

	prog, args := os.Args[2], os.Args[3:]
	cmd := exec.Command(prog, args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Command start error: %s", err)
		return
	}

	readLines(bufio.NewReader(stdout), "O: ", tgi)
	readLines(bufio.NewReader(stderr), "E: ", tgi)

	if err := cmd.Wait(); err != nil {
		fmt.Println("Command wait error:", err)
	}
}
