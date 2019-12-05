package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/hashicorp/go-uuid"
	"log"
	"os/exec"
	"strings"
)

type Worker struct {
	id     string
	driver Driver
	state WorkerState
}

type Driver uint

const (
	Raw Driver = iota
	Docker
)

type WorkerState uint

const (
	Sleeping WorkerState = iota
	Walking
	Working
)

func NewWorker(driver Driver) *Worker {
	id, _ := uuid.GenerateUUID()
	return &Worker{
		id: id,
		driver: driver,
		state: Sleeping,
	}
}

func (w *Worker) MatchAny(task *storage.Task) bool {
	log.Printf("matching task %d", task.ID)
	return true
}

func (w *Worker) Execute(task *storage.Task) ([]int, storage.TaskState){

}

func (w *Worker) ExecuteCmd(cmd string) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]
	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("%s", out)
}