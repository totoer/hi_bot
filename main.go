package main

import (
	"context"
	"fmt"

	"hi_bot/executor"
	"hi_bot/models"
	"hi_bot/proxies"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic("Config not readed")
	}

	var dbClient *mongo.Client
	dbClient, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("db_uri")))

	if err != nil {
		panic("Error")
	}

	if err := dbClient.Connect(context.Background()); err != nil {
		panic("Error")
	}

	router := executor.NewRouter()

	bots := models.FindAllBot(dbClient)
	for _, bot := range bots {
		router.Append(bot.Template, bot.Script)
	}

	// router.Append("!Bot (.+)", "test.lua")
	// router.Append("!Bot (Test)", "test.lua")
	// router.Append("!Bot (.*)", "test.lua")

	// c := make(chan *executor.Message)
	// r := make(chan []string)
	// testProxy := new(proxies.TestProxy)
	// go testProxy.Run(c, r)

	discordMessageChan := make(chan *executor.Message)
	discordResponseChan := make(chan []string)
	discordProxy := proxies.NewDiscordProxy(viper.GetString("discord_bot_token"))
	go discordProxy.Run(discordMessageChan, discordResponseChan)

	fmt.Println("Start listen proxies")
	for {
		select {
		case message := <-discordMessageChan:
			fmt.Println("Receive message from DiscordProxy")
			process(discordResponseChan, router, message)
		}
	}
}
