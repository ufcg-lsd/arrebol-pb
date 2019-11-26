package storage

import (
	"context"
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

func (s *Storage) SaveQueueSpec(q QueueSpec) (interface{}, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &q)
}

func (s *Storage) RetrieveQueue(queueId string) (*QueueSpec, error) {
	filter := bson.M{"_id": queueId}

	var q QueueSpec

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, er := context.WithTimeout(context.Background(), 10*time.Second)

	if er != nil {
		log.Println("Request timeout")
	}

	e := collection.FindOne(ctx, filter).Decode(&q)

	return &q, e
}

type QueueSpec struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}
