package worker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"log"
	"net/http"
)

const SignatureHeader string = "Signature"
const PublicKeyHeader string = "Public-Key"
const WrongBodyMsg string = "Maybe the body has a wrong shape"
const TokenKey string = "arrebol-worker-token"
type TokenResponse struct {
	ArrebolWorkerToken string
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		signature string
		encodedPublicKey string
		publicKey []byte
		_worker   *worker.Worker
		_token    token.Token
		queueId   uint
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
		WriteBadRequest(&w, WrongBodyMsg + ": " + err.Error())
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

	log.Println("Worker [" + _worker.ID + "] has been successfully joined")
	log.Println("Worker [" + _worker.ID + "] Token: " + _token.String())
	api.Write(w, http.StatusCreated, map[string]string{TokenKey: _token.String()})
}

func (a *WorkerApi) GetTask(w http.ResponseWriter, r *http.Request) {
	var (
		signature string
		err       error
		reqToken  string
	)

	if signature, err = GetHeader(r, SignatureHeader); err != nil {
		WriteBadRequest(&w, "The "+ SignatureHeader + " is not in the request header.")
	}

	params := mux.Vars(r)
	workerId := params["wid"]
	queueId := params["qid"]
	endpoint := "/v1/workers/"+workerId+"/queues/"+queueId +"/tasks"

	if ok, err := auth.CheckSignature([]byte(endpoint), []byte(signature), workerId); !ok || err != nil {
		api.Write(w, 401, api.ErrorResponse{
			Message: "The signature is not valid",
			Status:  401,
		})
	}

	if reqToken, err = GetHeader(r, TokenKey); err != nil {
		WriteBadRequest(&w, "The "+TokenKey+ " is not in the request header.")
	}

	parsedToken := token.Token(reqToken)
	if err = a.auth.Authorize(&parsedToken); err != nil {
		api.Write(w, 403, api.ErrorResponse{
			Message: "Authorization error. Invalid token.",
			Status:  403,
		})
	}

	


}

func GetHeader(r *http.Request, key string) (string, error) {
	log.Println("Getting header [" + key + "]")
	value := r.Header.Get(key)
	if value == "" {return "", errors.New("The header [" + key + "] was not found")}
	return value, nil
}

func WriteBadRequest(w *http.ResponseWriter, msg string) {
	log.Println("Bad Request: " + msg)
	api.Write(*w, http.StatusBadRequest, api.ErrorResponse{
		Message: msg,
		Status:  http.StatusBadRequest,
	})
}