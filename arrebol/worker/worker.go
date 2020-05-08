package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type WorkerSpec struct {
	vcpu  float32
	ram   float32
	image string
}

type TaskState uint8

const (
	TaskPending TaskState = iota
	TaskRunning
	TaskFinished
	TaskFailed
)

type Task struct {
	commands       []string
	reportInterval int64
	state          TaskState
	progress       int
}

type PBWorker struct {
	serverEndPoint string
	spec           *WorkerSpec
	address        string
	token          string
	id             string
	queueId        string
}

func (w *PBWorker) reportTask(task *Task, endingChannel chan int) {
	startTime := time.Now().Unix()
	for {
		select {
		case <-endingChannel:
			return
		default:
			currentTime := time.Now().Unix()
			if currentTime-startTime < task.reportInterval {
				time.Sleep(1 * time.Second)
				continue
			}
			url := w.serverEndPoint + "/workers/" + w.id + "/queues/" + w.queueId + "/tasks"
			requestBody, err := json.Marshal(task)

			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
			req.Header.Set("arrebol-worker-token", w.token)

			if err != nil {
				// handle error
				log.Fatal(err)
			}

			_, err = client.Do(req)
			if err != nil {
				// handle error
				log.Fatal(err)
			}

			startTime = currentTime
		}
	}

}

func (w *PBWorker) getTask() *Task {
	if w.queueId == "" {
		log.Println("The queue id has not been set yet")
		return nil
	}

	if w.token == "" {
		log.Println("The token has not been set yet")
		return nil
	}

	url := w.serverEndPoint + "/workers/" + w.id + "/queues/" + w.queueId
	requestBody, err := json.Marshal(&PBWorker{spec: &WorkerSpec{ram: w.spec.ram,
		image: w.spec.image, vcpu: w.spec.vcpu}, address: w.address, id: w.id})

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(requestBody))
	req.Header.Set("arrebol-worker-token", w.token)

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	defer resp.Body.Close()

	reqBody, err := ioutil.ReadAll(resp.Body)
	var task Task
	json.Unmarshal(reqBody, &task)

	return &task
}

func (w *PBWorker) subscribe() {
	url := w.serverEndPoint + "/workers"
	requestBody, err := json.Marshal(&PBWorker{spec: &WorkerSpec{ram: w.spec.ram,
		image: w.spec.image, vcpu: w.spec.vcpu}, address: w.address, id: w.id})

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	defer resp.Body.Close()

	reqBody, err := ioutil.ReadAll(resp.Body)
	var parsedBody map[string]string
	json.Unmarshal(reqBody, &parsedBody)

	w.token = parsedBody["arrebol-worker-token"]
	w.queueId = parsedBody["queue_id"]
}

func (w *PBWorker) execTask(task *Task) {
	endingChannel := make(chan int, 1)
	go w.reportTask(task, endingChannel)
	taskSize := len(task.commands)
	for i, command := range task.commands {
		exec.Command("/bin/bash", "-c", getPrjPath()+"arrebol/worker/bin/task-command-executor.sh -c " + command)
		task.progress = (i * 100) / taskSize
	}
	endingChannel <- 1
}

func getPrjPath() string{
	path_cmd := exec.Command("/bin/sh", "-c", "echo $GOPATH")
	path, _ := path_cmd.Output()
	path_str := strings.TrimSpace(string(path))
	return path_str + "/src/github.com/ufcg-lsd/arrebol-pb/"
}

func LoadWorker() PBWorker {
	path_cmd := exec.Command("/bin/sh", "-c", "echo $GOPATH")
	path, _ := path_cmd.Output()
	path_str := strings.TrimSpace(string(path))

	log.Println("Starting reading configuration process")
	// it must open the port and make all scripts executable
	file, err := os.Open(path_str + "/src/github.com/ufcg-lsd/arrebol-pb/arrebol/worker/worker-conf.json.example")
	defer file.Close()

	if err != nil {
		log.Println("Error on opening configuration file", err.Error())
	}

	decoder := json.NewDecoder(file)
	configuration := PBWorker{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("Error on decoding configuration file", err.Error())
	}

	return configuration
}
