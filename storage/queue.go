package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func (s *Storage) SaveQueue(q *Queue) (*mongo.InsertOneResult, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &q)
}

func (s *Storage) RetrieveQueue(queueId string) (*Queue, error) {
	var id primitive.ObjectID
	if queueId == DefaultQueueID {
		id, _ = getObjectIDFromDefault()
	} else {
		id, _ = primitive.ObjectIDFromHex(queueId)
	}

	filter := bson.M{"_id": id}

	var q Queue

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&q)

	if err != nil {
		log.Printf("%s not found in db", id.Hex())
	}

	return &q, err
}

func (s *Storage) RetrieveQueues() ([]*Queue, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	var queues []*Queue

	findOpt := options.Find()

	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOpt)

	if err != nil {
		log.Println("failed to retrieve queues")
	}

	if cursor != nil {
		for cursor.Next(context.TODO()) {
			var queue Queue
			err := cursor.Decode(&queue)
			if err != nil {
				log.Println("failed to decode queue")
			}
			queues = append(queues, &queue)
		}
		if err := cursor.Err(); err != nil {
			log.Printf("cursor error %v\n", err)
		}

		_ = cursor.Close(context.TODO())
	}
	return queues, nil
}

type ResourceState int8

const (
	Idle ResourceState = iota
	Allocated
	Busy
)

func (rs ResourceState) String() string {
	return [...]string{"Idle, Allocated, Busy"}[rs]
}

type Queue struct {
	ID   primitive.ObjectID `json:"ID,omitempty" bson:"_id,omitempty"`
	Name string             `json:"Name" bson:"name"`
	Jobs []Job              `json:"Jobs" bson:"jobs"`
	Nodes []ResourceNode 	 		`json:"Nodes" bson:"nodes"`
}

type ResourceNode struct {
	State ResourceState `json:"State" bson:"state"`
	Address string `json:"Address" bson:"address"`
}
