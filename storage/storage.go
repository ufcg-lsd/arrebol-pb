package storage

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"log"
)

type Storage struct {
	driver *gorm.DB
}

const dbDialect string =  "postgres"
var DB *Storage

func New(host string, port string, user string, dbname string, password string) *Storage {
	dbConfig := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	driver, err := gorm.Open(dbDialect, dbConfig)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = driver.DB().Ping()

	if err != nil {
		log.Fatalln(err.Error())
	}

	DB = &Storage{
		driver,
	}

	return DB
}

func (s *Storage) Setup() {
	s.CreateSchema()
	createDefaults(s)
}

func (s *Storage) Driver() *gorm.DB {
	return s.driver
}

func (s *Storage) SaveQueue(q *Queue) error {
	return s.driver.Save(&q).Error
}

func (s *Storage) RetrieveQueue(queueID uint) (*Queue, error) {
	var queue Queue
	err := s.driver.First(&queue, queueID).Error
	queue.Jobs, _ = s.RetrieveJobsByQueueID(queueID)
	queue.Workers, _ = s.RetrieveWorkersByQueueID(queueID)
	// TODO Retrieve resource nodes of queue
	return &queue, err
}

func (s *Storage) RetrieveQueues() ([]*Queue, error) {
	var queues []*Queue

	err := s.driver.Find(&queues).Error

	return queues, err
}

func (s *Storage) RetrieveJobs() ([]*Job, error) {
	var jobs []*Job

	err := s.driver.Find(&jobs).Error

	return jobs, err
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

func (s *Storage) RetrieveTask(taskId uint) (*Task, error){
	var task Task
	err := s.driver.First(&task, taskId).Error
	s.fillTask(&task)
	return &task, err
}

func (s *Storage) SaveJob(job *Job) error {
	return s.driver.Save(&job).Error
}

func (s *Storage) SetJobState(jobID uint, state JobState) {
	var job Job
	s.driver.First(&job, jobID)
	s.driver.Model(&job).Update("State", state)
}

func (s *Storage) SetTaskState(taskID uint, state TaskState) {
	var task Job
	s.driver.First(&task, taskID)
	s.driver.Model(&task).Update("State", state)
}

func (s *Storage) SaveTask(task *Task) error {
	return s.driver.Save(&task).Error
}

func (s *Storage) SaveCommand(command *Command) error {
	return s.driver.Save(&command).Error
}

func (s *Storage) RetrieveJobByQueue(jobID, queueId uint) (*Job, error) {
	var queue Queue
	var job Job

	err := s.driver.First(&queue, queueId).Related(&queue.Jobs).Error
	if queue.QueueHasJob(jobID) {
		err := s.driver.First(&job, jobID).Related(&job.Tasks).Error
		s.fillTasks(job.Tasks)
		return &job, err
	} else {
		err = errors.New(fmt.Sprintf("Job [%d] not found on queue [%d]", jobID, queueId))
	}
	return nil, err
}

func (s *Storage) RetrieveJobsByQueueID(queueID uint) ([]*Job, error) {
	var jobs []*Job

	log.Printf("Retrieving jobs of queue %d", queueID)
	err := s.driver.Where("queue_id = ?", queueID).Find(&jobs).Error

	for i, job := range jobs {
		s.driver.First(&job, job.ID).Related(&job.Tasks)
		s.fillTasks(job.Tasks)
		jobs[i] = job
	}

	return jobs, err
}

func (s *Storage) fillTasks(tasks []*Task) {
	for _, task := range tasks {
		s.fillTask(task)
	}
}

func (s *Storage) fillTask(task *Task) {
	db := s.driver.First(&task, task.ID)
	db.Related(&task.Commands)
	db.Related(&task.Metadata)
	db.Related(&task.Config)
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

func (s *Storage) RetrieveWorkersByQueueID(queueID uint) ([]*worker.Worker, error) {
	var workers []*worker.Worker

	log.Printf("Retrieving workers of queue %d", queueID)
	err := s.driver.Where("queue_id = ?", queueID).Find(&workers).Error

	return workers, err
}

func createDefaults(storage *Storage) {
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
