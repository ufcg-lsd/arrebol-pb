package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"

type Storage struct {
	client *mongo.Client
}

func New(client *mongo.Client, wait time.Duration) *Storage {
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	client, _ = mongo.Connect(ctx)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the storage")

	return &Storage{
		client,
	}
}

func (s *Storage) SaveQueue(q *Queue) (*mongo.InsertOneResult, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &q)
}

func (s *Storage) RetrieveQueue(queueId string) (*Queue, error) {
	id, _ := primitive.ObjectIDFromHex(queueId)
	filter := bson.M{"_id": id}

	var q Queue

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, er := context.WithTimeout(context.Background(), 10*time.Second)

	if er != nil {
		log.Println("Request timeout")
	}

	e := collection.FindOne(ctx, filter).Decode(&q)

	return &q, e
}

type Queue struct {
	ID   primitive.ObjectID `json:"ID,omitempty" bson:"_id,omitempty"`
	Name string `json:"Name" bson:"name"`
}
