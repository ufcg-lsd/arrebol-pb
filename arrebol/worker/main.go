package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

func setup(serverEndPoint string, workerId string) {
	cmd := exec.Command("/bin/sh", "-c",
		"$GOPATH/src/github.com/ufcg-lsd/arrebol-pb/arrebol/worker/bin/generate-ssh-key-pair.sh && cat "+workerId+"-key.pub")
	key, err := cmd.Output()

	if err != nil {
		log.Println("Error on generating worker key")
	}

	url := serverEndPoint + "/workers/publicKey"
	requestBody, err := json.Marshal(&map[string]string{"key": string(key)})

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
}

func main() {
	worker := LoadWorker()
	setup(worker.serverEndPoint, worker.id)
	worker.subscribe()
	task := worker.getTask()
	joinChan := make(chan interface{})

	go worker.execTask(task)

	<-joinChan
}
