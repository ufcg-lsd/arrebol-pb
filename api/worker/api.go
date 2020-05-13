package worker

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type WorkerApi struct {
	server  *http.Server
}

func New() *WorkerApi {
	return &WorkerApi{
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
	log.Print(mux.Vars(r))
	body := make(map[string]*rsa.PublicKey)
	json.NewDecoder(r.Body).Decode(&body)
	for k := range body {
		log.Println(k)
		log.Println(body[k])
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	log.Print(mux.Vars(r))
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	println(body["data"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
}

func (a *WorkerApi) GetTask(w http.ResponseWriter, r *http.Request) {
	log.Print(mux.Vars(r))
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	println(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

func (a *WorkerApi) ReportTask(w http.ResponseWriter, r *http.Request) {
	log.Print(mux.Vars(r))
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	println(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

