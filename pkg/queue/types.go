package queue

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkerNode struct {
	Address    string `json:"address" bson:"address"`
	WorkerPool int64  `json:"worker_pool" bson:"worker_pool"`
}

type Queue struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	WorkerNodes []WorkerNode       `json:"worker_nodes" bson:"worker_nodes"`
}
