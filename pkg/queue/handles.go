package queue

import (
	"context"
	"encoding/json"
	"github.com/emanueljoivo/arrebol/pkg/env"
	"github.com/emanueljoivo/arrebol/pkg/wrapper"
	"log"
	"net/http"
	"os"
	"time"
)

func CreateQueue(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var q Queue

	err := json.NewDecoder(r.Body).Decode(&q)

	if err != nil {
		log.Println("Error while process the request")
	}

	collection := wrapper.Client().Database(os.Getenv(env.DatabaseName)).Collection(os.Getenv(env.QueueCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, &q)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			log.Println(err.Error())
		}
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println(err.Error())
	}
}