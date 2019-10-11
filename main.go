package main

import (
	"fmt"
	"hi_bot/executor"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// var dbClient *mongo.Client

	// dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))

	// if err != nil {
	// 	panic("Error")
	// }

	// if err := dbClient.Connect(context.Background()); err != nil {
	// 	panic("Error")
	// }

	router := executor.NewRouter()

	// bots := models.FindAllBot(dbClient)
	// for _, bot := range bots {
	// 	router.Append(bot.Template, bot.Script)
	// }

	router.Append("!Bot (.+)", "test.lua")

	for {
		pipe := make(chan *executor.Message)
		go router.Run(pipe)
		message := executor.NewMessage("Test", "!Bot Test")
		pipe <- message
		select {
		case <-pipe:
			fmt.Println(message.Result)
			return
		}
	}
}
