package models

type WorkerNode struct {
	Address    string `json:"address" bson:"address"`
	WorkerPool int64  `json:"worker_pool" bson:"worker_pool"`
}