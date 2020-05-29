package worker

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"net/http"
)

const SignatureHeader string = "SIGNATURE";

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var (
		signature string
		worker worker.Worker
		queueId string
	)

	signature = r.Header.Get(SignatureHeader)

	if signature == "" {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "signature header was not found",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}

	data, err := json.Marshal(w)
	if err := a.auth.VerifySignature(worker.ID, data, []byte(signature)); err != nil {
		// TODO return
	}

	err = a.auth.VerifySignature(worker.ID, data, []byte(signature))
	if err != nil {
		//TODO return
	}

	queueId, err = a.manager.Join(worker)
	worker.QueueId = queueId

	if err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	token, err := a.auth.CreateToken(&worker)

	if err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	w.Header().Set("Token", (*token).String())
	api.Write(w, http.StatusOK, nil)
}

