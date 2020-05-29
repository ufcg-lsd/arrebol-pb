package worker

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/manager"
	"net/http"
)

const SignatureHeader string = "SIGNATURE";

func newAuthorization(signature string, message []byte) manager.Authorization {
	authorization := manager.Authorization{
		[]byte(signature),
		message,
	}
	return authorization
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get(SignatureHeader)
	//data, err := ioutil.ReadAll(r.Body)

	if signature == "" {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "signature header was not found",
			Status:  http.StatusBadRequest,
		})
		return
	}

	var worker worker.Worker
	//bytes.NewReader(data)
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}
	data, err := json.Marshal(worker)

	auth := newAuthorization(signature, data)
	token, err := a.workerManager.Join(auth, worker)

	if err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	w.Header().Set("Token", token.String())
	api.Write(w, http.StatusOK, nil)
}

