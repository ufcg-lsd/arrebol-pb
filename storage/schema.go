package storage

import (
	"github.com/jinzhu/gorm"
	"time"
)

func (s *Storage) CreateSchema() {
	s.driver.CreateTable(&Queue{})
	s.driver.CreateTable(&Job{})
	s.driver.CreateTable(&Task{})
	s.driver.CreateTable(&Command{})
}

func (s *Storage) UpdateSchema() {
	s.driver.AutoMigrate(&Queue{})
	s.driver.AutoMigrate(&Job{})
	s.driver.AutoMigrate(&Task{})
	s.driver.AutoMigrate(&Command{})
}

type Queue struct {
	ID        string         `json:"ID"`
	Name      string         `json:"Name"`
	Jobs      []Job          `json:"Jobs" gorm:"ForeignKey:QueueID"`
	Nodes     []ResourceNode `json:"Nodes" gorm:"ForeignKey:QueueID"`
}

type JobState uint8

const (
	JobPending JobState = iota
	JobRunning
	JobFinished
	JobFailed
)

func (js JobState) String() string {
	return [...]string{"Pending, Running, Failed", "Finished"}[js]
}

type Job struct {
	ID string `json:"ID"`
	QueueID string `json:"QueueID"`
	Label string   `json:"Label"`
	State JobState `json:"State"`
	Tasks []Task   `json:"Tasks" gorm:"ForeignKey:JobID"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	DeletedAt time.Time `json:"DeletedAt"`
}

type TaskState uint8

const (
	TaskPending TaskState = iota
	TaskRunning
	TaskFinished
	TaskFailed
)

func (ts TaskState) String() string {
	return [...]string{"Pending, Running, Failed", "Finished"}[ts]
}

type Task struct {
	gorm.Model
	JobID string `json:"JobID"`
	Config   []TaskConfig   `json:"Config" gorm:"ForeignKey:TaskID"`
	State    TaskState      `json:"State"`
	Metadata []TaskMetadata `json:"Metadata" gorm:"ForeignKey:TaskID"`
	Commands []Command      `json:"Commands" gorm:"ForeignKey:TaskID"`
}

type TaskConfig struct {
	Key   interface{}
	Value interface{}
}

type TaskMetadata struct {
	Key   interface{}
	Value interface{}
}

type CommandState uint8

const (
	CmdNotStarted CommandState = iota
	CmdRunning
	CmdFinished
	CmdFailed
)

func (cs CommandState) String() string {
	return [...]string{"NotStarted, Running, Failed", "Finished"}[cs]
}

type Command struct {
	gorm.Model
	TaskID int `json:"TaskID"`
	ExitCode   int8         `json:"Commands"`
	RawCommand string       `json:"RawCommand"`
	State      CommandState `json:"State"`
}

type ResourceState uint8

const (
	Idle ResourceState = iota
	Allocated
	Busy
)

func (rs ResourceState) String() string {
	return [...]string{"Idle, Allocated, Busy"}[rs]
}

type ResourceNode struct {
	gorm.Model
	State   ResourceState `json:"State"`
	Address string        `json:"Address"`
}
