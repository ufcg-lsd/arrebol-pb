package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/emanueljoivo/arrebol/pkg"
)

const CurrentVersion = "0.0.1"

func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pkg.Version{Tag: CurrentVersion}); err != nil {
		log.Println(err.Error())
	}
}
