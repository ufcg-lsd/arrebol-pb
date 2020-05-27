package worker

import (
	"bytes"
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWorkerApiAddWorker(t *testing.T) {
	_ = os.Setenv(crypto.KEYS_PATH, "../../test/keys")
	var workerSpec = WorkerSpec{
		ID:      "fake_worker",
		VCPU:    1.5,
		RAM:     1024,
		QueueId: "default",
	}

	workerSpecBytes, err := json.Marshal(workerSpec)
	if err != nil {
		t.Fatal(err)
	}

	var mockSignature []uint8
	log.Println(os.Getwd())
	if privKey, err := crypto.GetPrivateKey("../../test/keys/fake_worker"); err != nil {
		t.Fatal(err)
	} else {
		mockSignature, err = crypto.Sign(privKey, workerSpecBytes)
	}

	api := New()
	req, err := http.NewRequest("POST", "/v1/workers", bytes.NewBuffer(workerSpecBytes))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(SIGNATURE_HEADER, string(mockSignature))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.AddWorker)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"QueueId":"default","Token":"fake_worker#default#1.50#1024"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}