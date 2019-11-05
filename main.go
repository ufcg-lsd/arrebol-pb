package pkg

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/emanueljoivo/arrebol/pkg/handler"
)

const GetVersionEndpoint = "/version"
const CreateQueueEndpoint = "/queues"
const GetQueueEndpoint = "/queues/{id}"

func init() {
	log.Println("Starting Arrebol")

	ValidateEnv()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful_timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	SetUp(ctx)

	router := mux.NewRouter()

	router.HandleFunc(GetVersionEndpoint, handler.GetVersion).Methods("GET")
	router.HandleFunc(GetQueueEndpoint, handler.RetrieveQueue).Methods("GET")

	router.HandleFunc(CreateQueueEndpoint, handler.CreateQueue).Methods("POST")

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	log.Println("Service available")

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err.Error())
	}

	log.Println("Shutting down service")

	os.Exit(1)
}
