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
}

type Driver uint

const (
	Raw Driver = iota
	Docker
)

func NewWorker(driver Driver) *Worker {
	id, _ := uuid.GenerateUUID()
	return &Worker{
		id: id,
		driver: driver,
	}
}

func (w *Worker) MatchAny(task *storage.Task) bool {
	log.Printf("matching task %d", task.ID)
	return true
}

func (w *Worker) Execute(cmd string) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]
	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("%s", out)
}