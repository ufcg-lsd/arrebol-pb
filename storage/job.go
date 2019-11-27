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

func (s *Storage) SaveJob(job *Job, objId primitive.ObjectID) (*mongo.InsertOneResult, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &job)
}

func (s *Storage) RetrieveJob(jobId string) (*Job, error) {
	id, _ := primitive.ObjectIDFromHex(jobId)

	filter := bson.M{"_ID": id}

	var job Job

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&job)

	if err != nil {
		log.Printf("%s not found in db", id.Hex())
	}

	return &job, err
}

type State string

const (
	Pending  State = "Pending"
	Running  State = "Running"
	Finished State = "Finished"
	Failed   State = "Failed"
)

type Job struct {
	ID    primitive.ObjectID `json:"ID,omitempty" bson:"_id,omitempty"`
	Label string             `json:"Label" bson:"label"`
	State State              `json:"State" bson:"state"`
	Tasks []Task             `json:"Tasks" bson:"tasks"`
}

type Task struct {
	ID       string                 `json:"id,omitempty" bson:"_id,omitempty"`
	Config   map[string]interface{} `json:"Config" bson:"config"`
	Commands []string               `json:"Commands" bson:"commands"`
	State    State                  `json:"State" bson:"state"`
}
