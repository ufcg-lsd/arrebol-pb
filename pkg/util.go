package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"
const ServerPort = "SERVER_PORT"
const DefaultServerPort = "8080"

func ValidateEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file specified, searching for env variables...")
	}

	if _, exists := os.LookupEnv(DatabaseAddress); !exists {
		log.Fatal("No storage address on the env")
	} else if _, exists := os.LookupEnv(DatabaseName); !exists {
		log.Fatal("No storage name on the env")
	} else if _, exists := os.LookupEnv(ServerPort); !exists {
		log.Println("Setting default server port")
		err := os.Setenv(ServerPort, DefaultServerPort)
		if err != nil {
			log.Fatal("Maybe the port already in use")
		}
	} else {
		log.Println("Environment loaded with success")
	}
}
