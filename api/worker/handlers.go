package worker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"log"
	"net/http"
)

const SignatureHeader string = "Signature"
const PublicKeyHeader string = "Public-Key"
const WrongBodyMsg string = "Maybe the body has a wrong shape"

type TokenResponse struct {
	ArrebolWorkerToken string
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var (
		err              error
		signature        string
		encodedPublicKey string
		publicKey        []byte
		_worker          *worker.Worker
		_token           token.Token
		queueId          uint
	)

	if signature, err = GetHeader(r, SignatureHeader); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	if encodedPublicKey, err = GetHeader(r, PublicKeyHeader); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	if publicKey, err = base64.StdEncoding.DecodeString(encodedPublicKey); err != nil {
		WriteBadRequest(&w, err.Error())
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&_worker); err != nil {
		WriteBadRequest(&w, WrongBodyMsg+": "+err.Error())
		return
	}

	if _token, err = a.auth.Authenticate(string(publicKey), []byte(signature), _worker); err != nil {
		log.Println("Unauthorized: " + r.RemoteAddr + " - " + err.Error())
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	if err = a.auth.Authorize(&_token); err != nil {
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

	log.Println("Worker [" + _worker.ID.String() + "] has been successfully joined")
	log.Println("Worker [" + _worker.ID.String() + "] Token: " + _token.String())
	api.Write(w, http.StatusCreated, map[string]string{"arrebol-worker-token": _token.String()})
}

func GetHeader(r *http.Request, key string) (string, error) {
	log.Println("Getting header [" + key + "]")
	value := r.Header.Get(key)
	if value == "" {
		return "", errors.New("The header [" + key + "] was not found")
	}
	return value, nil
}

func WriteBadRequest(w *http.ResponseWriter, msg string) {
	log.Println("Bad Request: " + msg)
	api.Write(*w, http.StatusBadRequest, api.ErrorResponse{
		Message: msg,
		Status:  http.StatusBadRequest,
	})
}
