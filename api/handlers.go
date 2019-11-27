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

var (
	ProcReqErr = errors.New("error while trying to process response")
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

	queueId := params["queueId"]

	queue, err := a.storage.RetrieveQueue(queueId)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))

	} else {
		write(w, http.StatusOK, queue)
	}
}

func (a *API) CreateJob(w http.ResponseWriter, r *http.Request) {
	var jobSpec storage.JobSpec

	err := json.NewDecoder(r.Body).Decode(&jobSpec)

	if err != nil {
		log.Println(ProcReqErr)
	}

	id := primitive.NewObjectID()

	_, err = a.storage.SaveJob(&jobSpec, id)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"ID": "%s"}`, id.Hex())
	}
}

func (a *API) GetVersion(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, Version{Tag: VersionTag, Name: VersionName})
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
