package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-uuid"
	"log"
	"net/http"
	"time"
)

const VersionTag = "0.0.1"
const VersionName = "Havana"

type Version struct {
	Tag  string `json:"Tag"`
	Name string `json:"Name"`
}

type QueueResponse struct {
	ID           string `json:"ID"`
	Name         string `json:"Name"`
	PendingTasks uint   `json:"PendingTasks"`
	RunningTasks uint   `json:"RunningTasks"`
	Nodes        uint   `json:"Nodes"`
	Workers      uint   `json:"Workers"`
}

type JobResponse struct {
	ID        string     `json:"ID"`
	Label     string     `json:"Label"`
	State     string     `json:"State"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	Tasks     []TaskSpec `json:"Tasks"`
}

type ErrorResponse struct {
	Message string `json:"Message"`
	Status uint `json:"Status"`
}

type JobSpec struct {
	Label string     `json:"Label"`
	Tasks []TaskSpec `json:"Tasks"`
}

type TaskSpec struct {
	ID       string            `json:"ID"`
	Config   map[string]string `json:"Config"`
	Commands []string          `json:"Commands"`
	Metadata map[string]string `json:"Metadata"`
}

var (
	ProcReqErr   = errors.New("error while trying to process response")
	EncodeResErr = errors.New("error while trying encode response")
)

func (a *API) CreateQueue(w http.ResponseWriter, r *http.Request) {
	var q storage.Queue

	err := json.NewDecoder(r.Body).Decode(&q)

	if err != nil {
		log.Println(ProcReqErr)
	}

	q.ID, _ = uuid.GenerateUUID()

	a.storage.SaveQueue(&q)

	if err != nil {
		write(w, http.StatusBadRequest, ErrorResponse{
			Message: "Error while trying to save the new queue",
			Status:  http.StatusBadRequest,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"ID": "%s"}`, q.ID)
	}
}

func (a *API) RetrieveQueue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	queueID := params["qid"]

	queue := a.storage.RetrieveQueue(queueID)

	if queue.ID == "" {
		write(w, http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("Queue with id %s not found", queueID),
			Status:  http.StatusNotFound,
		})
	} else {
		pendingTasks := a.storage.RetrieveTasksByState(queue.ID, storage.TaskPending)
		runningTasks := a.storage.RetrieveTasksByState(queue.ID, storage.TaskRunning)
		response := responseFromQueue(queue, uint(len(pendingTasks)), uint(len(runningTasks)))

		write(w, http.StatusOK, &response)
	}
}

func (a *API) RetrieveQueues(w http.ResponseWriter, r *http.Request) {
	var response []*QueueResponse

	queues := a.storage.RetrieveQueues()

	for _, queue := range queues {
		pendingTasks := a.storage.RetrieveTasksByState(queue.ID, storage.TaskPending)
		runningTasks := a.storage.RetrieveTasksByState(queue.ID, storage.TaskRunning)
		curQueue := responseFromQueue(&queue, uint(len(pendingTasks)), uint(len(runningTasks)))
		response = append(response, curQueue)
	}
	write(w, http.StatusOK, response)
}

func (a *API) CreateJob(w http.ResponseWriter, r *http.Request) {
	var jobSpec JobSpec
	params := mux.Vars(r)

	queueID := params["qid"]

	err := json.NewDecoder(r.Body).Decode(&jobSpec)

	if err != nil {
		log.Println(ProcReqErr)
	}

	job := extractFromSpec(jobSpec)

	queue := a.storage.RetrieveQueue(queueID)
	queue.Jobs = append(queue.Jobs, *job)
	a.storage.SaveJob(job)
	a.storage.SaveQueue(queue)

	if err != nil {
		write(w, http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"ID": "%s"}`, job.ID)
	}
}

func (a *API) AddNode(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusAccepted, `{"Message": "no support yet"}`)
}

//func (a *API) RetrieveJobsByQueue(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//
//	queueId := params["qid"]
//
//	jobs, err := a.storage.RetrieveJobsByQueueID(queueId)
//
//	if err != nil {
//		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
//	} else {
//		write(w, http.StatusOK, jobs)
//	}
//}
//
//func (a *API) RetrieveJobByQueue(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//
//	queueId := params["qid"]
//	jobId := params["jid"]
//
//	job, err := a.storage.RetrieveJobByQueue(jobId, queueId)
//
//	if err != nil {
//		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
//	} else {
//		write(w, http.StatusOK, job)
//	}
//}
//

func (a *API) RetrieveNode(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusAccepted, `{"Message": "no support yet"}`)
}

func (a *API) RetrieveNodes(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusAccepted, `{"Message": "no support yet"}`)
}

func (a *API) GetVersion(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, Version{Tag: VersionTag, Name: VersionName})
}

func write(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Println(EncodeResErr)
	}
}

func responseFromQueue(queue *storage.Queue, pendingTasks uint, runningTasks uint) *QueueResponse {
	return &QueueResponse{
		ID:           queue.ID,
		Name:         queue.Name,
		PendingTasks: pendingTasks,
		RunningTasks: runningTasks,
		Nodes:        uint(len(queue.Nodes)),
	}
}

func extractFromSpec(spec JobSpec) *storage.Job {
	var tasks []storage.Task

	for _, taskSpec := range spec.Tasks {
		configs := extractConfigs(&taskSpec)
		metadata := extractMetadata(&taskSpec)
		commands := extractCommands(&taskSpec)

		tasks = append(tasks, storage.Task{
			Config:   configs,
			State:    storage.TaskPending,
			Metadata: metadata,
			Commands: commands,
		})
	}
	jobID, _ := uuid.GenerateUUID()
	now := time.Now()
	return &storage.Job{
		ID:   jobID,
		Label: spec.Label,
		State: storage.JobPending,
		Tasks: tasks,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func extractCommands(spec *TaskSpec) []storage.Command {
	var commands []storage.Command
	for _, cmd := range spec.Commands {
		commands = append(commands, storage.Command{
			ExitCode:   -1,
			RawCommand: cmd,
			State:      storage.CmdNotStarted,
		})
	}
	return commands
}

func extractMetadata(spec *TaskSpec) []storage.TaskMetadata {
	var metadata []storage.TaskMetadata
	for k, v := range spec.Config {
		metadata = append(metadata, storage.TaskMetadata{
			Key:   k,
			Value: v,
		})
	}
	return metadata
}

func extractConfigs(task *TaskSpec) []storage.TaskConfig {
	var configs []storage.TaskConfig
	for k, v := range task.Config {
		configs = append(configs, storage.TaskConfig{
			Key:   k,
			Value: v,
		})
	}

	return configs
}
