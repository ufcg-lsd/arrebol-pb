package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/emanueljoivo/arrebol/pkg/env"
	"github.com/emanueljoivo/arrebol/pkg/queue"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetVersion(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("0.0.1")
}

func CreateQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var queue queue.Queue

	_ = json.NewDecoder(r.Body).Decode(&queue)

	collection := client.Database(os.Getenv(env.DatabaseName)).Collection(os.Getenv(env.QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, &queue)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			// the blank field returns the number of bytes written
			log.Println(err.Error())
		}
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println(err.Error())
	}
}

func init() {
	log.Println("Starting Arrebol")

	flag.Parse()
	flag.Args()
	flag.Usage()

	env.ValidateEnv()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	clientOpt := options.Client().ApplyURI(os.Getenv(env.DatabaseAddress))
	client, _ = mongo.Connect(ctx, clientOpt)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the database")

	router := mux.NewRouter()

	router.HandleFunc("/version", GetVersion).Methods("GET")
	router.HandleFunc("/queues", CreateQueue).Methods("POST")

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	log.Println("Service available")

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err.Error())
	}

	log.Println("Shutting down service")

	os.Exit(1)
}
