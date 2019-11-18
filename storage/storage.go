package storage

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

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"

type Storage struct {
	client *mongo.Client
}

func New(ctx *context.Context, client *mongo.Client) *Storage {
	clientOpt := options.Client().ApplyURI(os.Getenv(DatabaseAddress))
	client, _ = mongo.Connect(*ctx, clientOpt)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the storage")

	return &Storage{
		client,
	}
}

func (s *Storage) SaveQueue(q Queue) (interface{}, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &q)
}

func (s *Storage) RetrieveQueue(queueId string) (*Queue, error) {

	qId, err := primitive.ObjectIDFromHex(queueId)

	if err != nil {
		log.Println("Queue id with wrong shape: " + queueId)
	}

	filter := bson.M{"_id": qId}

	var q Queue

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, er := context.WithTimeout(context.Background(), 10*time.Second)

	if er != nil {
		log.Println("Request timeout")
	}

	e := collection.FindOne(ctx, filter).Decode(&q)

	return &q, e
}
