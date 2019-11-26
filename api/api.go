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

	const GetVersionEndpoint = "/version"
	const CreateQueueEndpoint = "/queues"
	const GetQueueEndpoint = "/queues/{id}"

	router := mux.NewRouter()

	router.HandleFunc(GetVersionEndpoint, a.GetVersion).Methods("GET")
	router.HandleFunc(GetQueueEndpoint, a.RetrieveQueue).Methods("GET")

	router.HandleFunc(CreateQueueEndpoint, a.CreateQueue).Methods("POST")

	return router

}



