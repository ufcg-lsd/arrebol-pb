package driver

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/emanueljoivo/arrebol/helper"
	"github.com/emanueljoivo/arrebol/storage"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	TaskScriptExecutorFileName = "task-script-executor.sh"
	TaskScriptExecutorPath = "./resources/" + TaskScriptExecutorFileName
	RunTaskScriptCommandPattern = "/bin/bash %s -d -tsf=%s"
	PoolingPeriodTime = 2 * time.Second
)

type DockerDriver struct {
	Id      string
	Address string
	Cli     client.Client
}

func (d *DockerDriver) Execute(task storage.Task) error {
	task.State = storage.TaskRunning
	//TODO Receive docker image as a parameter
	config := helper.ContainerConfig{
		Name:   d.Id,
		Image:  "wesleymonte/simple-worker",
		Mounts: []mount.Mount{},
	}
	//TODO Handle errors
	d.initiate(config)
	d.send(task)
	d.run(strconv.Itoa(int(task.ID)))
	d.track(task)
	//TODO Check task state
	task.State = storage.TaskFinished
	return nil
}

func (d *DockerDriver) initiate(config helper.ContainerConfig) (err error) {
	exist, _ := helper.CheckImage(&d.Cli, config.Image)
	if !exist {
		//TODO Handle error and add log
		_, _ = helper.Pull(&d.Cli, config.Image)
	}
	cid, err := helper.CreateContainer(&d.Cli, config);
	if err == nil {
		err = helper.StartContainer(&d.Cli, cid)
		if err == nil {
			err = helper.Copy(&d.Cli, cid, TaskScriptExecutorPath, "/tmp/" + TaskScriptExecutorFileName)
		}
	}
	return err
}

func (d *DockerDriver) stop() error {
	err := helper.StopContainer(&d.Cli, d.Id)
	if err == nil {
		err = helper.RemoveContainer(&d.Cli, d.Id)
	}
	return err
}

func (d *DockerDriver) send(task storage.Task) error {
	taskScriptFileName := "task-id.ts"
	rawCmdsStr := task.GetRawCommands()
	// Check the maximum content size
	return helper.Write(&d.Cli, d.Id, rawCmdsStr, "/tmp/" + taskScriptFileName)
}

func (d *DockerDriver) run(taskId string) error {
	taskScriptFilePath := "/tmp/task-id.ts"
	cmd := fmt.Sprintf(RunTaskScriptCommandPattern, "/tmp/" + TaskScriptExecutorFileName, taskScriptFilePath)
	return helper.Exec(&d.Cli, d.Id, cmd)
}

func (d *DockerDriver) track(task storage.Task) error {
	i := 0
	for ; i < len(task.Commands);  {
		ec, _ := d.getExitCodes("task-id")
		i = d.syncCommands(task.Commands, ec, i)
		time.Sleep(PoolingPeriodTime)
	}
	return nil
}

func (d *DockerDriver) getExitCodes(taskId string) ([]int8, error) {
	ecFilePath := "/tmp/task-id" + ".ts.ec"
	dat, _ := helper.Cat(&d.Cli, d.Id, ecFilePath)
	dat = bytes.TrimFunc(dat, isNotUTFNumber)
	content := string(dat[:])
	log.Println("Content: " + content)
	exitCodesStr := strings.Split(content, "\r\n")
	log.Println("ExitCodes String Array: ", exitCodesStr)
	exitCodes := toIntArray(exitCodesStr)
	log.Println(exitCodes)
	return exitCodes, nil
}

func (d *DockerDriver) syncCommands(commands []*storage.Command, exitCodes []int8, startIndex int) int {
	i := startIndex
	for ; i < len(exitCodes) ; i++ {
		ec := exitCodes[i]
		if ec == 0 {
			commands[i].State = storage.CmdFinished
		} else {
			commands[i].State = storage.CmdFailed
		}
		commands[i].ExitCode = ec
	}
	if i < len(commands) {
		commands[i].State = storage.CmdRunning
	}
	return i
}

func toIntArray(strs []string) []int8 {
	ints := make([]int8, 0)
	for _, s := range strs {
		x, err := strconv.Atoi(s)
		if err == nil {
			ints = append(ints, int8(x))
		}
	}
	return ints
}

func isNotUTFNumber(r rune) bool {
	if r >= 48 && r <= 57 {
		return false
	}
	return true
}
