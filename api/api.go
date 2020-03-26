// Testing go-swagger generation
//
// The purpose of this application is to test go-swagger in a simple GET request.
//
//     Schemes: http
//     Host: localhost:8080
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Daniel<danielfs.ti@gmail.com>
//
//     Consumes:
//     - text/plain
//
//     Produces:
//     - text/plain
//
// swagger:meta
package api

import (
	"context"
	"github.com/emanueljoivo/arrebol/arrebol"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type HttpApi struct {
	storage *storage.Storage
	server  *http.Server
	arrebol *arrebol.Dispatcher
}

func New(storage *storage.Storage, arrebol *arrebol.Dispatcher) *HttpApi {
	return &HttpApi{
		storage: storage,
		arrebol: arrebol,
	}
}

func (a *HttpApi) Start(port string) error {
	a.server = &http.Server{
		Addr:    ":" + port,
		Handler: a.bootRouter(),
	}
	log.Println("Service available")
	return a.server.ListenAndServe()
}

func (a *HttpApi) Shutdown() error {
	return a.server.Shutdown(context.Background())
}

func (a *HttpApi) bootRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/version", a.GetVersion).Methods(http.MethodGet)

	router.HandleFunc("/v1/queues", a.CreateQueue).Methods(http.MethodPost)
	router.HandleFunc("/v1/queues", a.RetrieveQueues).Methods(http.MethodGet)
	router.HandleFunc("/v1/queues/{qid}", a.RetrieveQueue).Methods(http.MethodGet)

	router.HandleFunc("/v1/queues/{qid}/jobs", a.CreateJob).Methods(http.MethodPost)
	router.HandleFunc("/v1/queues/{qid}/jobs", a.RetrieveJobsByQueue).Methods(http.MethodGet)
	router.HandleFunc("/v1/queues/{qid}/jobs/{jid}", a.RetrieveJobByQueue).Methods(http.MethodGet)

	router.HandleFunc("/v1/queues/{qid}/nodes", a.AddNode).Methods(http.MethodPost)
	router.HandleFunc("/v1/queues/{qid}/nodes", a.RetrieveNodes).Methods(http.MethodGet)
	router.HandleFunc("/v1/queues/{qid}/nodes/{nid}", a.RetrieveNode).Methods(http.MethodGet)
	router.HandleFunc("/swagger.json", a.Swagger).Methods(http.MethodGet)

	return router
}
