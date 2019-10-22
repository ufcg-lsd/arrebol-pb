package queue

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkerNode struct {
	Address      string
	ResourcePool int64
}

type Queue struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	WorkerNodes []WorkerNode       `json:"worker_nodes bson:"worker_nodes"`
}
