package arrebol

import (
	"github.com/emanueljoivo/arrebol/arrebol/driver"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/hashicorp/go-uuid"
)

type Worker struct {
	id     string
	driver driver.Driver
	state  WorkerState
}

type WorkerState uint

const (
	Sleeping WorkerState = iota
	Working
	Busy
)

func NewWorker(driver driver.Driver) *Worker {
	id, _ := uuid.GenerateUUID()
	return &Worker{
		id: id,
		driver: driver,
		state: Sleeping,
	}
}

func (w *Worker) MatchAny(task *storage.Task) bool {
	return w.state == Sleeping
}

func (w *Worker) Execute(task *storage.Task) {
	w.state = Working
	task.State = storage.TaskRunning
	_ = storage.DB.SaveTask(task)
	w.driver.Execute(task)
	_ = storage.DB.SaveTask(task)
	w.state = Sleeping
}