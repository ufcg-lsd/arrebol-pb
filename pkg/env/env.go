package env

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
		log.Fatal("Error loading .env file")
	}


	if _, exists := os.LookupEnv(DatabaseAddress); !exists {
		log.Fatal("No wrapper address on the environment.")
	} else if _, exists := os.LookupEnv(DatabaseName); !exists {
		log.Fatal("No wrapper name on the environment.")
	} else {
		log.Println("Environment loaded with success.")
	}
}
