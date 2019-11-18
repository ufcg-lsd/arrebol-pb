package pkg

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func SetUp(ctx context.Context) {
	clientOpt := options.Client().ApplyURI(os.Getenv(DatabaseAddress))
	client, _ = mongo.Connect(ctx, clientOpt)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the storage")
}

func SaveQueue(q Queue) (interface{}, error) {
	collection := client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &q)
}

func RetrieveQueue(queueId string) (*Queue, error) {

	qId, err := primitive.ObjectIDFromHex(queueId)

	if err != nil {
		log.Println("Queue id with wrong shape: " + queueId)
	}

	filter := bson.M{"_id": qId}

	var q Queue

	collection := client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, er := context.WithTimeout(context.Background(), 10*time.Second)

	if er != nil {
		log.Println("Request timeout")
	}

	e := collection.FindOne(ctx, filter).Decode(&q)

	return &q, e
}
