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
const PublicKeyHeader string = "PUBLIC_KEY"
const WrongBodyMsg string = "Maybe the body has a wrong shape"

type TokenResponse struct {
	ArrebolWorkerToken string
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		signature string
		publicKey string
		_worker   *worker.Worker
		_token    token.Token
		queueId   uint
	)

	if signature, err = GetHeader(r, SignatureHeader); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	if publicKey, err = GetHeader(r, PublicKeyHeader); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&_worker); err != nil {
		WriteBadRequest(&w, WrongBodyMsg)
		return
	}

	if _token, err = a.auth.Authenticate(publicKey, []byte(signature), _worker); err != nil {
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

	if _token, err = _token.SetPayloadField("QueueId", queueId); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	api.Write(w, http.StatusOK, TokenResponse{_token.String()})
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