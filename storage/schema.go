package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
)

func (s *Storage) DropTablesIfExist() *gorm.DB {
	return s.driver.DropTableIfExists(&Command{}, &TaskConfig{}, &TaskMetadata{},
		&Task{}, &Job{}, &ResourceNode{}, &Queue{}, &worker.Worker{})
}

func (s *Storage) CreateTables() {
	var tables = map[string]interface{}{
		"commands":       &Command{},
		"task_configs":   &TaskConfig{},
		"task_metadata":  &TaskMetadata{},
		"tasks":          &Task{},
		"jobs":           &Job{},
		"resource_nodes": &ResourceNode{},
		"queues":         &Queue{},
		"workers":        &worker.Worker{},
	}

	for _, v := range tables {
		err, _ := s.CreateTable(v)

		if err != nil {

		}
	}
}

func (s *Storage) CreateTable(t interface{}) (error, string) {
	clone := s.driver.CreateTable(t)

	if clone.Error != nil {
		var errMsg = fmt.Sprintf("Table %+v already exists", t)
		return errors.New(errMsg), errMsg
	} else {
		var successMsg = fmt.Sprintf("Table %+v correctly created", t)
		return nil, successMsg
	}
}

func (s *Storage) AutoMigrate() {
	s.driver.AutoMigrate(&Command{}, &TaskConfig{}, &TaskMetadata{},
		&Task{}, &Job{}, &ResourceNode{}, &Queue{})
}

func (s *Storage) ConfigureSchema() {
	s.Driver().Model(
		&Command{}).AddForeignKey(
		"task_id", "tasks(id)", "CASCADE", "CASCADE").Model(
		&TaskMetadata{}).AddForeignKey(
		"task_id", "tasks(id)", "CASCADE", "CASCADE").Model(
		&TaskConfig{}).AddForeignKey(
		"task_id", "tasks(id)", "CASCADE", "CASCADE").Model(
		&Task{}).AddForeignKey(
		"job_id", "jobs(id)", "CASCADE", "CASCADE").Model(
		&ResourceNode{}).AddForeignKey(
		"queue_id", "queues(id)", "CASCADE", "CASCADE").Model(
		&Job{}).AddForeignKey(
		"queue_id", "queues(id)", "CASCADE", "CASCADE").Model(
		&worker.Worker{}).AddForeignKey("queue_id", "queues(id)", "CASCADE", "CASCADE")
}

func (s *Storage) CreateSchema() {
	s.DropTablesIfExist()
	s.CreateTables()
	s.AutoMigrate()
	s.ConfigureSchema()
}

// swagger:model Queue
type Queue struct {
	gorm.Model
	Name    string           `json:"Name"`
	Jobs    []*Job           `json:"Jobs" gorm:"ForeignKey:QueueID"`
	Workers []*worker.Worker `json:"Workers" gorm:"ForeignKey:QueueID"`
	Nodes   []*ResourceNode  `json:"Nodes" gorm:"ForeignKey:QueueID"`
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
	return [...]string{"Queued", "Running", "Finished", "Failed"}[js]
}

type Job struct {
	gorm.Model
	QueueID uint     `json:"QueueID"`
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
	return [...]string{"Pending", "Running", "Finished", "Failed"}[ts]
}

type Task struct {
	gorm.Model
	JobID    uint           `json:"JobID"`
	State    TaskState      `json:"State"`
	Config   []TaskConfig   `json:"Config" gorm:"ForeignKey:TaskID"`
	Metadata []TaskMetadata `json:"Metadata" gorm:"ForeignKey:TaskID"`
	Commands []*Command     `json:"Commands" gorm:"ForeignKey:TaskID"`
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
	return [...]string{"NotStarted", "Running", "Finished", "Failed"}[cs]
}

type Command struct {
	gorm.Model
	TaskID     uint         `json:"TaskID"`
	ExitCode   int8         `json:"ExitCode"`
	RawCommand string       `json:"RawCommand"`
	State      CommandState `json:"State"`
}
