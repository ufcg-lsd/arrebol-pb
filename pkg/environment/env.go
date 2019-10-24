package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const QueueCollection = "QUEUE_COLLECTION"

func ValidateEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .environment file")
	}


	if _, exists := os.LookupEnv(DatabaseAddress); !exists {
		log.Fatal("No storage address on the environment.")
	} else if _, exists := os.LookupEnv(DatabaseName); !exists {
		log.Fatal("No storage name on the environment.")
	} else {
		log.Println("Environment loaded with success.")
	}
}
