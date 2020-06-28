package storage

import (
	"errors"
	"fmt"
	"log"
)

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

func (s *Storage) SetJobState(jobID uint, state JobState) {
	var job Job
	s.driver.First(&job, jobID)
	s.driver.Model(&job).Update("State", state)
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
