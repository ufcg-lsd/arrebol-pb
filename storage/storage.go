package storage

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

const DatabaseAddr = "DATABASE_ADDRESS"
const DatabasePort = "DATABASE_PORT"
const DatabaseName = "DATABASE_NAME"
const DatabasePassword = "DATABASE_PASSWORD"
const DatabaseUser = "DATABASE_USER"

type Storage struct {
	driver *gorm.DB
}

var (
	SaveErr = errors.New("error while trying to save document")
)

func New() *Storage {
	dbAddr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv(DatabaseAddr), os.Getenv(DatabasePort), os.Getenv(DatabaseUser), os.Getenv(DatabaseName), os.Getenv(DatabasePassword))
	driver, err := gorm.Open("postgres", dbAddr)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = driver.DB().Ping()

	if err != nil {
		log.Fatalln(err.Error())
	}

	driver.LogMode(true)

	storage := &Storage{
		driver,
	}

	return storage
}

func (s *Storage) SetUp() {
	s.CreateSchema()
	CreateDefault(s)
}

func (s *Storage) Driver() *gorm.DB {
	return s.driver
}

func (s *Storage) SaveQueue(q *Queue) error {
	return s.driver.Save(&q).Error
}

func (s *Storage) RetrieveQueue(queueId uint) (*Queue, error) {
	var queue Queue
	log.Println(fmt.Sprintf("retrieving queue %d", queueId))
	err := s.driver.First(&queue, queueId).Error
	log.Println(queue)
	return &queue, err
}

func (s *Storage) RetrieveQueues() []*Queue {
	var queues []*Queue

	s.driver.Find(&queues)

	return queues
}

func (s *Storage) RetrieveTasksByState(queueID uint, state TaskState) []*Task {
	var tasksPending []*Task
	queue, _ := s.RetrieveQueue(queueID)

	for _, job := range queue.Jobs {
		for _, task := range job.Tasks {
			if task.State == state {
				tasksPending = append(tasksPending, task)
			}
		}
	}

	return tasksPending
}

func (s *Storage) SaveJob(job *Job) {
	s.driver.Create(&job)
	s.driver.Save(&job)
}

//func (s *Storage) RetrieveJobByQueue(jobId string, queueId string) *Job {
//	var job Job
//
//	s.driver.Where()
//
//}
//
//func (s *Storage) RetrieveTasksByJob(jobID string) []*Task {
//
//}
//
//func (s *Storage) RetrieveCommandsByTask(taskID string) []*Command {
//
//}
//
//func (s *Storage) RetrieveJobsByQueueID(queueID string) ([]*Job, error) {
//
//}
//
//func generateObjID(queueID string) *primitive.ObjectID {
//	var id primitive.ObjectID
//	if queueID == DefaultQueueID {
//		id, _ = getObjectIDFromDefault()
//	} else {
//		id, _ = primitive.ObjectIDFromHex(queueID)
//	}
//	return &id
//}
//
//func getObjectIDFromDefault() (primitive.ObjectID, error) {
//	src := []byte(DefaultQueueID)
//
//	dst := make([]byte, hex.EncodedLen(len(src)))
//	hex.Encode(dst, src)
//
//	return primitive.ObjectIDFromHex(string(dst))
//}
//

func CreateDefault(storage *Storage) {
	q := &Queue{
		Name: "Default",
	}

	var queue Queue
	if err := storage.driver.Where("id = ?", q.ID).First(&queue).Error; err != nil {
		log.Println(err.Error())
		storage.SaveQueue(q)
	} else {
		log.Println("Default queue already exists")
	}
}
