package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/api/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const ServerPort = "SERVER_PORT"
	const DefaultServerPort = "5000"

	var wait time.Duration
	flag.DurationVar(&wait, "graceful_timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	apiPort := flag.String(ServerPort, DefaultServerPort, "Service port")

	flag.Parse()

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	s := storage.New(os.Getenv("DATABASE_ADDRESS"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	s.Setup()
	defer s.Driver().Close()

	var jobDispatcher = service.NewDispatcher(s)
	go jobDispatcher.Start()

	a := api.New(s, jobDispatcher)

	// Shutdown gracefully
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		log.Println("Shutting down service")

		if err := a.Shutdown(); err != nil {
			log.Printf("Failed to shutdown the server: %v", err)
		}
	}()

	go startWorkerApi(s)

	if err := a.Start(*apiPort); err != nil {
		log.Fatal(err.Error())
	}
}

func startWorkerApi(storage *storage.Storage) {
	const WorkerApiPort = "8000"

	workerApi := worker.New(storage)
	err := workerApi.Start(WorkerApiPort)

	if err != nil {
		log.Fatal(err.Error())
	}
}
