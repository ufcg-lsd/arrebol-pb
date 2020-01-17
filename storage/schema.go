package storage

import (
	"github.com/jinzhu/gorm"
)

func (s *Storage) CreateSchema() {
	s.driver.DropTable(&Command{}, &TaskConfig{}, &TaskMetadata{},
		&Task{}, &Job{}, &ResourceNode{}, &Queue{})

	s.driver.CreateTable(&Command{}, &TaskConfig{}, &TaskMetadata{},
		&Task{}, &Job{}, &ResourceNode{}, &Queue{})

	s.driver.AutoMigrate(&Command{}, &TaskConfig{}, &TaskMetadata{},
		&Task{}, &Job{}, &ResourceNode{}, &Queue{})

	s.Driver().Model(&TaskMetadata{}).AddForeignKey("task_id", "tasks(id)", "CASCADE", "CASCADE")
	s.Driver().Model(&TaskConfig{}).AddForeignKey("task_id", "tasks(id)", "CASCADE", "CASCADE")
	s.Driver().Model(&Command{}).AddForeignKey("task_id", "tasks(id)", "CASCADE", "CASCADE")

	s.Driver().Model(&Task{}).AddForeignKey("job_id", "jobs(id)", "CASCADE", "CASCADE")

	s.Driver().Model(&ResourceNode{}).AddForeignKey("queue_id", "queues(id)", "CASCADE", "CASCADE")
	s.Driver().Model(&Job{}).AddForeignKey("queue_id", "queues(id)", "CASCADE", "CASCADE")
}

type Queue struct {
	gorm.Model
	Name  string          `json:"Name"`
	Jobs  []*Job          `json:"Jobs" gorm:"ForeignKey:QueueID"`
	Nodes []*ResourceNode `json:"Nodes" gorm:"ForeignKey:QueueID"`
}

func (q Queue) contains(jobId uint) bool {
	for _, job := range q.Jobs {
		if job.ID == jobId {
			return true
		}
	}
	return false
}

type ResourceState uint8

const (
	Idle ResourceState = iota
	Allocated
)

func (rs ResourceState) String() string {
	return [...]string{"Idle", "Allocated"}[rs]
}

type ResourceNode struct {
	gorm.Model
	QueueID uint          `json:"QueueID"`
	State   ResourceState `json:"State"`
	Address string        `json:"Address"`
}

type JobState uint8

const (
	JobQueued JobState = iota
	JobRunning
	JobFinished
	JobFailed
)

func (js JobState) String() string {
	return [...]string{"Queued", "Running", "Failed", "Finished"}[js]
}

type Job struct {
	gorm.Model
	QueueID uint   `json:"QueueID"`
	Label   string   `json:"Label"`
	State   JobState `json:"State"`
	Tasks   []*Task  `json:"Tasks" gorm:"ForeignKey:JobID"`
}

type TaskState uint8

const (
	TaskPending TaskState = iota
	TaskRunning
	TaskFinished
	TaskFailed
)

func (ts TaskState) String() string {
	return [...]string{"Pending", "Running", "Failed", "Finished"}[ts]
}

type Task struct {
	gorm.Model
	JobID    uint           `json:"JobID"`
	State    TaskState      `json:"State"`
	Config   []TaskConfig   `json:"Config" gorm:"ForeignKey:TaskID"`
	Metadata []TaskMetadata `json:"Metadata" gorm:"ForeignKey:TaskID"`
	Commands []*Command      `json:"Commands" gorm:"ForeignKey:TaskID"`
}

type TaskConfig struct {
	gorm.Model
	TaskID uint
	Key    string
	Value  string
}

type TaskMetadata struct {
	gorm.Model
	TaskID uint
	Key    string
	Value  string
}

type CommandState uint8

const (
	CmdNotStarted CommandState = iota
	CmdRunning
	CmdFinished
	CmdFailed
)

func (cs CommandState) String() string {
	return [...]string{"NotStarted", "Running", "Failed", "Finished"}[cs]
}

type Command struct {
	gorm.Model
	TaskID     uint         `json:"TaskID"`
	ExitCode   int8         `json:"Commands"`
	RawCommand string       `json:"RawCommand"`
	State      CommandState `json:"State"`
}
