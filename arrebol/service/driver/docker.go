package driver

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/docker"
	"github.com/ufcg-lsd/arrebol-pb/storage"
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
	DockerImagePropertyKey = "docker_image"
	DefaultWorkerDockerImage = "wesleymonte/simple-worker"
)

const (
	PullImageErrorMsg       		string = "Error while pulling docker image [%s]"
	CreateContainerErrorMsg 		string = "Error while creating container [%s]"
	StartContainerErrorMsg  		string = "Error while starting container [%s]"
	CopyTaskScriptErrorMsg  		string = "Error while driver [%s] copying task script executor [%s]"
	StopContainerErrorMsg   		string = "Error while stopping container [%s]"
	RemoveContainerErrorMsg 		string = "Error while removing container [%s]"
	SendTaskScriptFileErrorMsg 		string = "Error while send task script file [%s]"
	RunTaskScriptExecutorErrorMsg 	string = "Error while running the " + TaskScriptExecutorFileName
	GettingExitCodesErrorMsg 		string = "Error while getting content of exit codes file [%s]"
	TrackTaskErrorMsg				string = "Error while track task execution [%s]"
)


type DockerDriver struct {
	Id      string
	Cli     client.Client
}

func (d *DockerDriver) Execute(task *storage.Task) error {
	image, err := task.GetConfig(DockerImagePropertyKey)
	if err == nil {
		image = DefaultWorkerDockerImage
	}
	config := docker.ContainerConfig{
		Name:   d.Id,
		Image:  image,
		Mounts: []mount.Mount{},
	}
	if err = d.initiate(config); err != nil {
		return err
	}
	if err = d.send(task); err != nil {
		return err
	}
	if err = d.run(strconv.Itoa(int(task.ID))); err != nil {
		return err
	}
	if err = d.track(task); err != nil {
		return err
	}
	if err = d.stop(); err != nil {
		return err
	}
	task.State = storage.TaskFinished
	return nil
}

func (d *DockerDriver) initiate(config docker.ContainerConfig) error {
	exist, err := docker.CheckImage(&d.Cli, config.Image)
	if !exist {
		if _, err = docker.Pull(&d.Cli, config.Image); err != nil {
			return errors.Wrapf(err, PullImageErrorMsg, config.Image)
		}
	}
	cid, err := docker.CreateContainer(&d.Cli, config);
	if err != nil {
		return errors.Wrapf(err, CreateContainerErrorMsg, config.Name)
	}
	err = docker.StartContainer(&d.Cli, cid)
	if err != nil {
		return errors.Wrapf(err, StartContainerErrorMsg, config.Name)
	}
	err = docker.Copy(&d.Cli, cid, TaskScriptExecutorPath, "/tmp/" + TaskScriptExecutorFileName)
	if err != nil {
		return errors.Wrapf(err, CopyTaskScriptErrorMsg, d.Id, TaskScriptExecutorFileName)
	}
	return err
}

func (d *DockerDriver) stop() error {
	err := docker.StopContainer(&d.Cli, d.Id)
	if err != nil {
		return errors.Wrapf(err, StopContainerErrorMsg, d.Id)
	}
	err = docker.RemoveContainer(&d.Cli, d.Id)
	if err != nil {
		return errors.Wrapf(err, RemoveContainerErrorMsg, d.Id)
	}
	return err
}

func (d *DockerDriver) send(task *storage.Task) error {
	taskScriptFileName := "task-id.ts"
	rawCmdsStr := task.GetRawCommands()
	err := docker.Write(&d.Cli, d.Id, rawCmdsStr, "/tmp/" + taskScriptFileName)
	if err != nil {
		err = errors.Wrapf(err, SendTaskScriptFileErrorMsg, taskScriptFileName)
	}
	return err
}

func (d *DockerDriver) run(taskId string) error {
	taskScriptFilePath := "/tmp/task-id.ts"
	cmd := fmt.Sprintf(RunTaskScriptCommandPattern, "/tmp/" + TaskScriptExecutorFileName, taskScriptFilePath)
	err := docker.Exec(&d.Cli, d.Id, cmd)
	if err != nil {
		err = errors.Wrap(err, RunTaskScriptExecutorErrorMsg)
	}
	return err
}

func (d *DockerDriver) track(task *storage.Task) error {
	i := 0
	for ; i < len(task.Commands);  {
		ec, err := d.getExitCodes("task-id")
		if err != nil {
			return errors.Wrapf(err, TrackTaskErrorMsg, "task-id")
		}
		i = d.syncCommands(task.Commands, ec, i)
		time.Sleep(PoolingPeriodTime)
	}
	return nil
}

func (d *DockerDriver) getExitCodes(taskId string) ([]int8, error) {
	ecFilePath := "/tmp/task-id" + ".ts.ec"
	dat, err := docker.Cat(&d.Cli, d.Id, ecFilePath)
	if err != nil {
		err = errors.Wrapf(err, GettingExitCodesErrorMsg, ecFilePath)
		return nil, err
	}
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
		storage.DB.SaveCommand(commands[i])
	}
	if i < len(commands) {
		commands[i].State = storage.CmdRunning
		storage.DB.SaveCommand(commands[i])
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
