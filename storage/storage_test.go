package storage

import "testing"

func TestSaveQueue(t *testing.T) {
	s := OpenDriver()

	s.Setup()
	var tasks []*Task
	tasks = append(tasks, &Task{
		JobID:    0,
		State:    0,
		Config:   nil,
		Metadata: nil,
		Commands: nil,
	},
	)

	err := s.SaveJob(&Job{
		QueueID: 1,
		Label:   "Some Label",
		State:   0,
		Tasks: tasks,
	})

	if err != nil {
		t.Error(err)
	}
}