package wrapper

import (
	"context"
	"github.com/emanueljoivo/arrebol/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var client *mongo.Client

func Client() *mongo.Client {
	return client
}

func SetUp(ctx context.Context) {
	clientOpt := options.Client().ApplyURI(os.Getenv(env.DatabaseAddress))
	client, _ = mongo.Connect(ctx, clientOpt)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the wrapper")
}