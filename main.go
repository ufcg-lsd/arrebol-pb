package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emanueljoivo/arrebol/api"
	"github.com/emanueljoivo/arrebol/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	const ServerPort = "SERVER_PORT"
	const DefaultServerPort = "8080"

	var wait time.Duration
	flag.DurationVar(&wait, "graceful_timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	apiPort := flag.String(ServerPort, DefaultServerPort, "Service port")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	m := mongo.NewClient(&mongo.Op)
	s := storage.New(ctx)
	a := api.New(s)

	// Shutdown gracefully
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		log.Println("Shutting down service")

		if err := a.Shutdown(&ctx); err != nil {
			log.Fatal(err.Error())
		}
	}()

	if err := a.Start(*apiPort); err != nil {
		log.Fatal(err.Error())
	}

}
