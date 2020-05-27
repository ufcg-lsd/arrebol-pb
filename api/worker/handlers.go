package worker

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"io/ioutil"
	"net/http"
	"strings"
)

const SIGNATURE_HEADER string = "SIGNATURE";

type WorkerSpec struct {
	ID      string  `json:"ID"`
	VCPU    float32 `json:"vCPU"`
	RAM     uint32  `json:"RAM"` //Megabytes
	QueueId string  `json:"QueueId, omitempty"`
}

type Token struct {
	WorkerId string `json:"WorkerId"`
	QueueId string  `json:"QueueId"`
	VCPU	float32 `json:"vCPU"`
	RAM     uint32  `json:"RAM"`
}

func (t Token) getRaw() string {
	const separator string = "#"
	builder := strings.Builder{}
	builder.WriteString(t.WorkerId)
	builder.WriteString(separator)
	builder.WriteString(t.QueueId)
	builder.WriteString(separator)
	builder.WriteString(fmt.Sprintf("%.2f", t.VCPU))
	builder.WriteString(separator)
	builder.WriteString(fmt.Sprint(t.RAM))
	return builder.String()
}

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	err := a.verifySignature(r)
	if err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	var workerSpec WorkerSpec
	_ = json.NewDecoder(r.Body).Decode(&workerSpec)
	queueId := a.selectQueue(workerSpec)
	token := a.generateToken(queueId, workerSpec)
	api.Write(w, http.StatusOK, struct {
		QueueId string `json:"QueueId"`
		Token   string `json:"Token"`
	}{queueId, token.getRaw()})
}

func (a *WorkerApi) verifySignature(r *http.Request) (err error) {
	signature := r.Header.Get(SIGNATURE_HEADER)
	body, _ := r.GetBody()

	if signature == "" {
		return errors.New("request signature was not found")
	}

	var workerSpec WorkerSpec
	if err = json.NewDecoder(body).Decode(&workerSpec); err != nil {
		return errors.New("Maybe the body has a wrong shape")
	}

	var publicKey *rsa.PublicKey
	var message []byte

	body, _ = r.GetBody()
	if publicKey, err = crypto.GetWorkerPublicKey(workerSpec.ID); err != nil {return}
	if message, err = ioutil.ReadAll(body); err != nil {return}

	return crypto.Verify(publicKey, message, []byte(signature))
}

func (a *WorkerApi) selectQueue(workerSpec WorkerSpec) string {
	return ""
}

func (a *WorkerApi) generateToken(queueId string, spec WorkerSpec) Token {
	return Token{}
}

