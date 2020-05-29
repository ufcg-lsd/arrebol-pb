package worker

import (
	"github.com/gorilla/mux"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/manager"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"net/http"
)

type WorkerApi struct {
	server  *http.Server
	workerManager *manager.Manager
	storage *storage.Storage
}

func New(storage *storage.Storage) *WorkerApi {
	return &WorkerApi{
		storage:       storage,
		workerManager: manager.NewManager(),
	}
}

func (a *WorkerApi) Start(port string) error {
	a.server= &http.Server{
		Addr:    ":" + port,
		Handler: a.bootRouter(),
	}
	log.Println("Starting worker api")
	return a.server.ListenAndServe()
}

func (a *WorkerApi) bootRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/workers", a.AddWorker).Methods(http.MethodPost)
	router.HandleFunc("/v1/workers/publicKey", a.AddPublicKey).Methods(http.MethodPost)
	router.HandleFunc("/v1/workers/{wid}/queues/{qid}/tasks", a.GetTask).Methods(http.MethodGet)
	router.HandleFunc("/v1/workers/{wid}/queues/{qid}/tasks", a.ReportTask).Methods(http.MethodPut)

	return router
}

func (a *WorkerApi) AddPublicKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
}

func (a *WorkerApi) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

func (a *WorkerApi) ReportTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

