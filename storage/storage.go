package storage

import (
	"context"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"
const DefaultQueueID = "default-uuid"

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

	storage := &Storage{
		client,
	}

	CreateDefault(storage)

	return storage
}

func getObjectIDFromDefault() (primitive.ObjectID, error) {
	src := []byte(DefaultQueueID)

	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return primitive.ObjectIDFromHex(string(dst))
}

func CreateDefault(storage *Storage) {
	id, err := getObjectIDFromDefault()

	if err != nil {
		log.Println(err.Error())
	}

	_, err = storage.RetrieveQueue(id.Hex())

	if err != nil {
		q := &Queue{
			Name: "Default",
			ID:   id,
		}
		_, err = storage.SaveQueue(q)

		if err != nil {
			log.Fatalln("error while trying create the default queue")
		}
	} else {
		log.Println("Queue default already exists")
	}
}
