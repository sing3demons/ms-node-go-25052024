package store

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStore(ctx context.Context) *mongo.Client {
	url := os.Getenv("MONGO_URI")

	loggerOptions := options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url).SetLoggerOptions(loggerOptions))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	return client
}
