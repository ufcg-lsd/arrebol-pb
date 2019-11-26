package api

import (
	"encoding/json"
	"errors"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-uuid"
	"log"
	"net/http"
)

const VersionTag = "0.0.1"
const VersionName = "Havana"

type Version struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

var (
	ProcReqErr = errors.New("error while trying to process response")
	EncodeResErr = errors.New("error while trying encode response")
)

func (a *API) CreateQueue(w http.ResponseWriter, r *http.Request) {
	var q storage.QueueSpec

	err := json.NewDecoder(r.Body).Decode(&q)

	if err != nil {
		log.Println(ProcReqErr)
	}

	q.ID, _ = uuid.GenerateUUID()

	res, err := a.storage.SaveQueueSpec(q)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err))
	} else {
		write(w, http.StatusCreated, res)
	}
}

func (a *API) RetrieveQueue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	queueId := params["id"]

	queue, err := a.storage.RetrieveQueue(queueId)

	if err != nil {
		write(w, http.StatusInternalServerError, notOkResponse(err))
	} else {
		write(w, http.StatusOK, queue)
	}
}

func (a *API) GetVersion(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, Version{Tag: VersionTag, Name: VersionName})
}

func notOkResponse(err error) []byte {
	return []byte(`{"message": ` + err.Error())
}

func write(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Println(EncodeResErr)
	}
}
