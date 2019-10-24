package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Queue struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	WorkerNodes []WorkerNode       `json:"worker_nodes" bson:"worker_nodes"`
}
