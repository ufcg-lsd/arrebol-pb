package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/emanueljoivo/arrebol/pkg/env"
	"github.com/emanueljoivo/arrebol/pkg/queue"
	"github.com/emanueljoivo/arrebol/pkg/wrapper"
)

func init() {
	log.Println("Starting Arrebol")

	env.ValidateEnv()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	wrapper.SetUp(ctx)

	router := mux.NewRouter()

	router.HandleFunc(env.VersionEndpoint, env.GetVersion).Methods("GET")
	router.HandleFunc(env.QueueEndpoint, queue.CreateQueue).Methods("POST")

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
