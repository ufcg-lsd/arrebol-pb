package env

import (
	"log"
	"os"
)

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"

func ValidateEnv() {
	if _, exists := os.LookupEnv(DatabaseAddress); !exists {
		log.Fatal("No database address on the environment.")
	} else if _, exists := os.LookupEnv(DatabaseName); !exists {
		log.Fatal("No database name on the environment.")
	} else {
		log.Println("Environment loaded with success.")
	}
}
