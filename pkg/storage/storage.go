package storage

import (
	"context"
	"github.com/emanueljoivo/arrebol/pkg/environment"
	"github.com/emanueljoivo/arrebol/pkg/queues"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client

func SetUp(ctx context.Context) {
	clientOpt := options.Client().ApplyURI(os.Getenv(environment.DatabaseAddress))
	client, _ = mongo.Connect(ctx, clientOpt)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the storage")
}

func SaveQueue(q queues.Queue) (interface{}, error) {
	collection := client.Database(os.Getenv(environment.DatabaseName)).Collection(os.Getenv(environment.QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &q)
}

func RetrieveQueue(queueId string) (*queues.Queue, error) {

	qId, err := primitive.ObjectIDFromHex(queueId)

	if err != nil {
		log.Println("Queue id with wrong shape: " + queueId)
	}

	filter := bson.M{"_id": qId}

	var q queues.Queue

	collection := client.Database(os.Getenv(environment.DatabaseName)).Collection(os.Getenv(environment.QueueCollection))

	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

	e := collection.FindOne(ctx, filter).Decode(&q)

	return &q, e
}