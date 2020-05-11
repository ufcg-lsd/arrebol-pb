package driver

import (
	"github.com/ufcg-lsd/arrebol/storage"
	"log"
	"os/exec"
	"strings"
)

const (
	FailExitCode = 1
	SuccessExitCode = 0
)

type Driver interface {
	Execute(t *storage.Task) error
}

type RawDriver struct {}

func (r *RawDriver) Execute(task *storage.Task) error {
	flawed := false
	for _, cmd := range task.Commands {
		r.execute(cmd)
		if cmd.State == storage.CmdFailed {
			flawed = true
		}
	}

	if flawed {
		task.State = storage.TaskFailed
	} else {
		task.State = storage.TaskFinished
	}
	return nil
}

func (r *RawDriver) execute(cmd *storage.Command) {
	cmd.State = storage.CmdRunning
	_ = storage.DB.SaveCommand(cmd)
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
	_ = storage.DB.SaveCommand(cmd)
}