package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

func (s *Storage) SaveQueue(q *Queue) (*mongo.InsertOneResult, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(JobCollection))

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

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(JobCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&q)

	if err != nil {
		log.Printf("%s not found in db", id.Hex())
	}

	return &q, err
}

func (s *Storage) SaveJob(jobSpec *JobSpec, objId primitive.ObjectID) (*mongo.InsertOneResult, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(JobCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &jobSpec)
}

type State string

const (
	Queued   State = "Queued"
	Running  State = "Running"
	Finished State = "Finished"
	Failed   State = "Failed"
)

type Job struct {
	Label string  `json:"Label" bson:"label"`
	State State   `json:"State" bson:"state"`
	Spec  JobSpec `json:"Spec" bson:"spec"`
}

type JobSpec struct {

}

type Queue struct {
	ID   primitive.ObjectID `json:"ID,omitempty" bson:"_id,omitempty"`
	Name string             `json:"Name" bson:"name"`
}
