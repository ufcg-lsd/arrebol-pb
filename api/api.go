package api

import (
	"context"
	"log"
	"net/http"

	"github.com/emanueljoivo/arrebol/handler"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
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

func (a *API) Shutdown(context *context.Context) error {
	return a.server.Shutdown(context.Background())
}

func (a *API) bootRouter() *mux.Router {

	const GetVersionEndpoint = "/version"
	const CreateQueueEndpoint = "/queues"
	const GetQueueEndpoint = "/queues/{id}"

	router := mux.NewRouter()

	router.HandleFunc(GetVersionEndpoint, handler.GetVersion).Methods("GET")
	router.HandleFunc(GetQueueEndpoint, handler.RetrieveQueue).Methods("GET")

	router.HandleFunc(CreateQueueEndpoint, handler.CreateQueue).Methods("POST")

	return router

}
