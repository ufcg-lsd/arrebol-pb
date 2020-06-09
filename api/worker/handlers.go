package worker

import (
	"encoding/json"
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"net/http"
)

const SignatureHeader string = "SIGNATURE";
const WrongBodyMsg string = "Maybe the body has a wrong shape"

type TokenResponse struct {
	ArrebolWorkerToken string
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		signature string
	)

	if signature, err = GetHeader(r, SignatureHeader); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	var (
		_worker *worker.Worker
		t       token.Token
		queueId uint
	)

	if err = json.NewDecoder(r.Body).Decode(&_worker); err != nil {
		WriteBadRequest(&w, WrongBodyMsg)
		return
	}

	if t, err = a.auth.Authenticate([]byte(signature), _worker); err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	if queueId, err = a.manager.Join(*_worker); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	if t, err = t.SetPayloadField("QueueId", queueId); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	api.Write(w, http.StatusOK, TokenResponse{t.String()})
}

func GetHeader(r *http.Request, key string) (string, error) {
	value := r.Header.Get(key)
	if value == "" {return "", errors.New("The header [" + key + "] was not found")}
	return value, nil
}

func WriteBadRequest(w *http.ResponseWriter, msg string) {
	api.Write(*w, http.StatusBadRequest, api.ErrorResponse{
		Message: msg,
		Status:  http.StatusBadRequest,
	})
}