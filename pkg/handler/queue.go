package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/emanueljoivo/arrebol/pkg"
	"github.com/gorilla/mux"
)

func CreateQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var q pkg.Queue

	err := json.NewDecoder(r.Body).Decode(&q)

	if err != nil {
		log.Println("Error while process the request")
	}

	res, err := pkg.SaveQueue(q)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			log.Println(err.Error())
		}
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err.Error())
	}
}

func RetrieveQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	queueId := params["id"]

	queue, err := pkg.RetrieveQueue(queueId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			log.Println(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(queue); err != nil {
			log.Println(err.Error())
		}
	}
}
