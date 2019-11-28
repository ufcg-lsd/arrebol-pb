package storage

import (
	"context"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"
const DefaultQueueID = "default-uuid"

type Storage struct {
	client *mongo.Client
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
	CreatedAt time.Time 	 `json:"CreatedAt" bson:"created_at"`
	UpdatedAt time.Time		 `json:"UpdatedAt" bson:"updated_at"`
}

type Task struct {
	ID       string            `json:"ID,omitempty" bson:"_id,omitempty"`
	Config   map[string]string `json:"Config" bson:"config"`
	Commands []string          `json:"Commands" bson:"commands"`
	State    State             `json:"State" bson:"state"`
	Metadata map[string]string `json:"Metadata" bson:"Metadata"`
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
	ID    primitive.ObjectID `json:"ID,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"Name" bson:"name"`
	Jobs  []Job              `json:"Jobs" bson:"jobs"`
	Nodes []ResourceNode     `json:"Nodes" bson:"nodes"`
}

type ResourceNode struct {
	State   ResourceState `json:"State" bson:"state"`
	Address string        `json:"Address" bson:"address"`
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
			Jobs: []Job{},
			Nodes: []ResourceNode{},
		}
		_, err = storage.SaveQueue(q)

		if err != nil {
			log.Fatalln("error while trying create the default queue")
		}
	} else {
		log.Println("Queue default already exists")
	}
}

func (s *Storage) SaveQueue(q *Queue) (*mongo.InsertOneResult, error) {
	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q.Jobs = []Job{}
	q.Nodes = []ResourceNode{}

	return collection.InsertOne(ctx, &q)
}

func (s *Storage) RetrieveQueue(queueId string) (*Queue, error) {
	id := generateObjID(queueId)

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

func (s *Storage) EnqueueJob(job *Job, queueId string) {
	coll := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))
	job.UpdatedAt = time.Now()
	id := generateObjID(queueId)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$addToSet", bson.M{
		"jobs": &job,
	}}}

	res, err := coll.UpdateOne(context.Background(), filter, update)

	log.Printf("error %v", err)

	log.Printf("updated %v", res)
}

func (s *Storage) RetrieveJobByQueue(jobId string, queueId string) (*Job, error) {

	collection := s.client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jid, _ := primitive.ObjectIDFromHex(jobId)
	qid := generateObjID(queueId)
	filter := bson.M{ "_id" : qid }

	var queue Queue

	err := collection.FindOne(ctx, filter).Decode(&queue)

	if err != nil {
		log.Printf("%s not found in db", jid.Hex())
	}
	var job Job
	jobs := queue.Jobs
	for _, value := range jobs {
		if value.ID == jid {
			job = value
			break
		}
	}

	return &job, err
}
func generateObjID(queueID string) *primitive.ObjectID {
	var id primitive.ObjectID
	if queueID == DefaultQueueID {
		id, _ = getObjectIDFromDefault()
	} else {
		id, _ = primitive.ObjectIDFromHex(queueID)
	}
	return &id
}
