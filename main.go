package main

import (
	"flag"
	"io"
	"log"
	"os"

	"hi_bot/executor"
	"hi_bot/proxies"
	"hi_bot/web"

	"github.com/spf13/viper"
)

type Proxy interface {
	Run(chan *executor.Message, chan []string)
}

func process(r chan []string, router *executor.Router, message *executor.Message) {
	pipe := make(chan *executor.Message)
	go router.Run(pipe)
	pipe <- message
	select {
	case <-pipe:
		r <- message.Result
	}
}

func main() {
	var err error

	configFile := flag.String("config", "config", "Config file")
	flag.Parse()

	viper.SetConfigName(*configFile)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		panic("Config not readed")
	}

	logFile, err := os.OpenFile(viper.GetString("logfile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Logfile not opening!")
	}
	defer logFile.Close()

	logMWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logMWriter)

	router := executor.NewRouter()

	go web.Run()

	discordMessageChan := make(chan *executor.Message)
	discordResponseChan := make(chan []string)
	discordQuitChan := make(chan int)
	discordProxy := proxies.NewDiscordProxy(viper.GetString("discord_bot_token"))
	go discordProxy.Run(discordMessageChan, discordResponseChan, discordQuitChan)

	log.Println("Start listen proxies")
	for {
		select {
		case message := <-discordMessageChan:
			log.Println("Receive message from DiscordProxy")
			process(discordResponseChan, router, message)
		}
	}

	defer func() {
		discordQuitChan <- 0
	}()
}
