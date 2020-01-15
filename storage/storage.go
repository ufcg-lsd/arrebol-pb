package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type Storage struct {
	driver *gorm.DB
}

func NewDB(host string, port string, user string, dbname string, password string) *Storage {
	dbAddr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	driver, err := gorm.Open("postgres", dbAddr)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = driver.DB().Ping()

	if err != nil {
		log.Fatalln(err.Error())
	}

	// driver.LogMode(true)

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

func (s *Storage) RetrieveQueue(queueID uint) (*Queue, error) {
	var queue Queue
	log.Println(fmt.Sprintf("Retrieving queue %d", queueID))
	err := s.driver.First(&queue, queueID).Error
	log.Println(queue)
	return &queue, err
}

func (s *Storage) RetrieveQueues() ([]*Queue, error) {
	var queues []*Queue

	err := s.driver.Find(&queues).Error

	return queues, err
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

func (s *Storage) SaveJob(job *Job) error {
	return s.driver.Save(&job).Error
}

func (s *Storage) RetrieveJobByQueue(jobID, queueId uint) (*Job, error) {
	var job Job
	log.Println(fmt.Sprintf("Retrieving job %d of queue %d", jobID, queueId))
	err := s.driver.First(&job, jobID).Error

	return &job, err
}

func (s *Storage) RetrieveJobsByQueueID(queueID uint) ([]Job, error) {
	var jobs []Job

	log.Printf("Retrieving jobs of queue %d", queueID)
	err := s.driver.Where("queue_id = ?", queueID).Find(&jobs).Error

	return jobs, err
}

func (s *Storage) GetDefaultQueue() (*Queue, error) {
	var queue Queue
	const QIDDefault = 1
	if err := s.driver.Where("id = ?", QIDDefault).First(&queue).Error; err == nil {
		s.driver.First(&queue, 1)
	} else {
		return nil, err
	}
	return &queue, nil
}

func CreateDefault(storage *Storage) {
	q := &Queue{
		Name: "Default",
	}

	queue, err := storage.GetDefaultQueue()

	if queue == nil {
		err = storage.SaveQueue(q)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println("Default queue already exists")
	}
}
