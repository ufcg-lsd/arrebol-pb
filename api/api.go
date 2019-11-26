package api

import (
	"context"
	"github.com/emanueljoivo/arrebol/storage"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	a.createDefault()
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

func (a *API) createDefault() {
	id, e := primitive.ObjectIDFromHex("000000000000")
	if e != nil {
		log.Println(e.Error())
	}
	q := &storage.Queue{
		Name: "Default",
		ID: id,
	}
	_, err := a.storage.SaveQueue(q)

	if err != nil {
		log.Fatalln("error while trying create the default queue")
	}
}

func (a *API) bootRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/version", a.GetVersion).Methods(http.MethodGet)

	router.HandleFunc("/v1/queues/{queueId}", a.RetrieveQueue).Methods(http.MethodGet)
	router.HandleFunc("/v1/queues", a.CreateQueue).Methods(http.MethodPost)

	return router
}



