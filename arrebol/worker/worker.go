package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
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

func reportReq(w *PBWorker, task *Task) {
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
}

func (w *PBWorker) reportTask(task *Task, endingChannel chan int) {
	startTime := time.Now().Unix()
	for {
		select {
		case <-endingChannel:
			task.state = TaskFinished
			reportReq(w, task)
			return
		default:
			currentTime := time.Now().Unix()
			if currentTime-startTime < task.reportInterval {
				time.Sleep(1 * time.Second)
				continue
			}

			reportReq(w, task)

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

func getPrivateKey(id string) *rsa.PrivateKey{
	readPrivKey, err := ioutil.ReadFile(getPrjPath()+"arrebol/worker/bin/"+id+".priv")
	if err != nil {
		log.Fatal("The private key is not where it should be")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(readPrivKey)
	if err != nil {
		log.Fatal("Error on parsing private key")
	}

	return rsaPrivateKey
}

func signMessage(privateKey *rsa.PrivateKey, message []byte) ([]byte, []byte){
	messageHash := sha256.New()
	_, err := messageHash.Write(message)
	if err != nil {
		panic(err)
	}
	msgHashSum := messageHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	return signature, msgHashSum
}

func getPublicKey(id string) *rsa.PublicKey{
	readPubKey, err := ioutil.ReadFile(getPrjPath()+"arrebol/worker/bin/"+id+".pub")
	if err != nil {
		log.Fatal("The private key is not where it should be")
	}

	rsaPubKey, err := x509.ParsePKCS1PublicKey(readPubKey)
	if err != nil {
		log.Fatal("Error on parsing private key")
	}

	return rsaPubKey
}

func verifySignature(key rsa.PublicKey, hash []byte, signature []byte) bool {
	err := rsa.VerifyPSS(&key, crypto.SHA256, hash, signature, nil)
	if err != nil {
		return false
	}
	return true
}

func (w *PBWorker) subscribe() {
	requestBody, err := json.Marshal(&PBWorker{spec: &WorkerSpec{ram: w.spec.ram,
		image: w.spec.image, vcpu: w.spec.vcpu}, address: w.address, id: w.id})

	data, hashSum := signMessage(getPrivateKey(w.id), requestBody)

	payload, _ := json.Marshal(&map[string][]byte{"data": data, "hashSum": hashSum})

	url := w.serverEndPoint + "/workers"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
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
