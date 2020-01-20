package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/hashicorp/go-uuid"
	"log"
	"os/exec"
	"strings"
)

const (
	FailExitCode = 1
	SuccessExitCode = 0
)

type Worker struct {
	id     	string
	driver 	Driver
	state 	WorkerState
}

type Driver uint

const (
	Raw Driver = iota
	Docker
)

type WorkerState uint

const (
	Sleeping WorkerState = iota
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
	return w.state == Sleeping
}

func (w *Worker) Execute(task *storage.Task) (storage.TaskState){
	w.state = Working
	task.State = storage.TaskRunning
	flawed := false
	for _, cmd := range task.Commands {
		w.ExecuteCmd(cmd)
		if cmd.State == storage.CmdFailed {
			flawed = true
		}
	}

	if flawed {
		task.State = storage.TaskFailed
	} else {
		task.State = storage.TaskFinished
	}
	w.state = Sleeping
	return task.State
}

func (w *Worker) ExecuteCmd(cmd *storage.Command) {
	cmd.State = storage.CmdRunning
	storage.DB.SaveCommand(cmd)
	cmdStr := cmd.RawCommand
	parts := strings.Fields(cmdStr)
	head := parts[0]
	parts = parts[1:]
	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		log.Printf("%s", err)
		cmd.State = storage.CmdFailed
		cmd.ExitCode = FailExitCode
	} else {
		log.Printf("%s", out)
		cmd.State = storage.CmdFinished
		cmd.ExitCode = SuccessExitCode
	}
	storage.DB.SaveCommand(cmd)
}