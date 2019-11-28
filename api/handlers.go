package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
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

	id := primitive.NewObjectID()
	q.ID = id

	_, err = a.storage.SaveQueue(&q)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"ID": "%s"}`, id.Hex())
	}
}

func (a *API) RetrieveQueue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	queueId := params["qid"]

	queue, err := a.storage.RetrieveQueue(queueId)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		response := responseFromQueue(queue)
		write(w, http.StatusOK, response)
	}
}

func (a *API) RetrieveQueues(w http.ResponseWriter, r *http.Request) {
	var response []*QueueResponse

	queues, err := a.storage.RetrieveQueues()

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		for i := 0; i < len(queues); i++ {
			curQueue := responseFromQueue(queues[i])
			response = append(response, curQueue)
		}
		write(w, http.StatusOK, response)
	}
}

func (a *API) CreateJob(w http.ResponseWriter, r *http.Request) {
	var jobSpec JobSpec
	params := mux.Vars(r)

	queueID := params["qid"]

	err := json.NewDecoder(r.Body).Decode(&jobSpec)

	if err != nil {
		log.Println(ProcReqErr)
	}

	jobID := primitive.NewObjectID()

	job, tasks, cmds := extractFromSpec(jobSpec, jobID, queueID)

	a.storage.SaveJob(job)
	a.storage.SaveTasks(tasks)
	a.storage.SaveCommands(cmds)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"ID": "%s"}`, jobID.Hex())
	}
}

func (a *API) RetrieveJobsByQueue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	queueId := params["qid"]

	jobs, err := a.storage.RetrieveJobsByQueueID(queueId)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		write(w, http.StatusOK, jobs)
	}
}

func (a *API) RetrieveJobByQueue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	queueId := params["qid"]
	jobId := params["jid"]

	job, err := a.storage.RetrieveJobByQueue(jobId, queueId)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		write(w, http.StatusOK, job)
	}
}

func (a *API) AddNode(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusAccepted, `{"message": "no support yet"}`)
}

func (a *API) RetrieveNode(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusAccepted, `{"message": "no support yet"}`)
}

func (a *API) RetrieveNodes(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusAccepted, `{"message": "no support yet"}`)
}

func (a *API) GetVersion(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, Version{Tag: VersionTag, Name: VersionName})
}

func responseFromQueue(queue *storage.Queue) *QueueResponse {
	var pendingTasks uint
	var runningTasks uint

	return &QueueResponse{
		ID:           queue.ID.Hex(),
		Name:         queue.Name,
		PendingTasks: pendingTasks,
		RunningTasks: runningTasks,
		Nodes:        0,
	}
}

func extractFromSpec(spec JobSpec, jobID primitive.ObjectID, queueID string) (*storage.Job,
	[]*storage.Task, []*storage.Command) {

	job := &storage.Job{
		ID:        jobID,
		Label:     spec.Label,
		State:     storage.JobPending,
		CreatedAt: jobID.Timestamp(),
		UpdatedAt: jobID.Timestamp(),
	}
	job.QueueID = queueID

	var tss []*storage.Task
	var cmds []*storage.Command

	for _, task := range spec.Tasks {
		tss = append(tss, &storage.Task{
			ID:       task.ID,
			JobID:    jobID.Hex(),
			Config:   task.Config,
			State:    storage.TaskPending,
			Metadata: task.Metadata,
		})
		for _, cmd := range task.Commands {
			cmds = append(cmds, &storage.Command{
				ExitCode:   -1,
				RawCommand: cmd,
				TaskID:     task.ID,
				State:      storage.CmdNotStarted,
			})
		}
	}
	return job, tss, cmds
}

func notOkResponse(err string) []byte {
	return []byte(`{"Message": ` + err)
}

func write(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Println(EncodeResErr)
	}
}
