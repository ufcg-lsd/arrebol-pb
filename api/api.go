package api

import (
	"context"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type API struct {
	storage *storage.Storage
	server  *http.Server
}

func New(storage *storage.Storage) *API {
	return &API{
		storage: storage,
	}
}

func (a *API) Start(port string) error {
	a.server = &http.Server{
		Addr:    ":" + port,
		Handler: a.bootRouter(),
	}
	log.Println("Service available")
	return a.server.ListenAndServe()
}

func (a *API) Shutdown() error {
	return a.server.Shutdown(context.Background())
}

func (a *API) bootRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/version", a.GetVersion).Methods(http.MethodGet)

	router.HandleFunc("/v1/queues", a.CreateQueue).Methods(http.MethodPost)
	router.HandleFunc("/v1/queues", a.RetrieveQueues).Methods(http.MethodGet)
	router.HandleFunc("/v1/queues/{qid}", a.RetrieveQueue).Methods(http.MethodGet)

	// router.HandleFunc("/v1/queues/{qid}/jobs", a.RetrieveJobsByQueue).Methods(http.MethodGet)
	// router.HandleFunc("/v1/queues/{qid}/jobs", a.CreateJob).Methods(http.MethodPost)
	// router.HandleFunc("/v1/queues/{qid}/jobs/{jid}", a.RetrieveJob).Methods(http.MethodGet)

	return router
}
