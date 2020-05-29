package worker

import (
	"bytes"
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/whitelist"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWorkerApiAddWorker(t *testing.T) {
	_ = os.Setenv(key.KeysPath, "../../test/keys")
	_ = os.Setenv(whitelist.WhiteListPath, "../../test/whitelist/whitelist")
	var worker = worker.Worker{
		ID:      "fake_worker",
		VCPU:    1.5,
		RAM:     1024,
		QueueId: "default",
	}

	data, err := json.Marshal(worker)
	if err != nil {
		t.Fatal(err)
	}

	var signature []uint8
	log.Println(os.Getwd())
	if privKey, err := crypto.GetPrivateKey("../../test/keys/fake_worker"); err != nil {
		t.Fatal(err)
	} else {
		signature, err = crypto.Sign(privKey, data)
	}

	api := New(nil)
	req, err := http.NewRequest("POST", "/v1/workers", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(SignatureHeader, string(signature))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.AddWorker)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `null`
	expectedToken := "fake_worker#default#1.50#1024"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if rr.Header().Get("Token") != expectedToken {
		t.Errorf("handler returned unexpected token: got %v want %v",
			rr.Body.String(), expected)
	}
}