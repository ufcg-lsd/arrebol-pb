package storage

import (
	"testing"
)

var storage = NewDB("localhost", "5432", "postgres", "postgres", "postgres")

func TestSaveQueue(t *testing.T) {
	storage.SetUp()
	var tasks []*Task
	tasks = append(tasks, &Task{
			JobID:    0,
			State:    0,
			Config:   nil,
			Metadata: nil,
			Commands: nil,
		},
	)

	err := storage.SaveJob(&Job{
		QueueID: 1,
		Label:   "Some Label",
		State:   0,
		Tasks: tasks,
	})

	if err != nil {
		t.Error(err)
	}
}