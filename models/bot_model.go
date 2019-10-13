package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BotModel struct {
	Template string
	Name     string
	Script   string
}

func FindAllBot(dbClient *mongo.Client) []BotModel {
	var result []BotModel
	var cursor *mongo.Cursor
	ctx := context.Background()

	collection := dbClient.Database("hi_bot").Collection("bots")
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		panic("Error")
	}

	for cursor.Next(ctx) {
		item := BotModel{}
		cursor.Decode(&item)
		result = append(result, item)
	}

	cursor.Close(ctx)
	return result
}

func (b *BotModel) Save(dbClient *mongo.Client) {}

func (b *BotModel) Delete(dbClient *mongo.Client) {}
