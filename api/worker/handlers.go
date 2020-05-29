package worker

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"net/http"
)

const SignatureHeader string = "SIGNATURE";

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var err error
	signature := r.Header.Get(SignatureHeader)

	if signature == "" {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "signature header was not found",
			Status:  http.StatusBadRequest,
		})
		return
	}

	var worker *worker.Worker
	if err = json.NewDecoder(r.Body).Decode(&worker); err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if err = a.auth.Validate([]byte(signature), worker); err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	var queueId string
	if queueId, err = a.manager.Join(*worker); err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}
	worker.QueueId = queueId
	token, err := a.auth.Authenticate(worker)

	if err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}

	w.Header().Set("Token", (*token).String())
	api.Write(w, http.StatusOK, nil)
}

