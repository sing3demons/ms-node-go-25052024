package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStore(ctx context.Context) *mongo.Client {
	url := "mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-replica-set"

	loggerOptions := options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url).SetLoggerOptions(loggerOptions))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	return client
}
