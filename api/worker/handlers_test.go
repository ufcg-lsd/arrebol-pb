package worker

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/allowlist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const FakeWorkerId = "fake_worker"

func OpenDriver() *storage.Storage {
	s := storage.New(os.Getenv("DATABASE_ADDRESS"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	s.Setup()
	return s
}

func CloseDriver(s *storage.Storage, t *testing.T) {
	err := s.Driver().Close()

	if err != nil {
		t.Fail()
	}
}

func GenerateWorkerSignature(body interface{}, workerId string) (signature []uint8, err error) {
	var data []byte
	if data, err = json.Marshal(body); err == nil {
		var privKey *rsa.PrivateKey
		if privKey, err = crypto.GetPrivateKey(os.Getenv(key.KeysPath) + "/" + workerId); err == nil {
			signature, err = crypto.Sign(privKey, data)
		}
	}
	return
}

func TestWorkerApiAddWorker(t *testing.T) {
	_ = godotenv.Load("../../test/.env")
	_ = os.Setenv(key.KeysPath, "../../test/keys")
	_ = os.Setenv(allowlist.ListFilePath, "../../test/allowlist/allowlist")
	s := OpenDriver()
	defer CloseDriver(s, t)

	worker := worker.Worker{
		ID:      FakeWorkerId,
		VCPU:    1.5,
		RAM:     1024,
		QueueID: 1,
	}
	data, err := json.Marshal(worker)
	CheckError(t, err)

	signature, err := GenerateWorkerSignature(worker, FakeWorkerId)
	CheckError(t, err)

	publicKey, err := ioutil.ReadFile("../../test/keys/fake.pub")
	CheckError(t, err)
	log.Println("PublicKey: " + string(publicKey))

	encodedPubKey := base64.StdEncoding.EncodeToString(publicKey)
	log.Println("Encoded PublicKey: " + encodedPubKey)

	api := New(s)
	req, err := http.NewRequest("POST", "/v1/workers", bytes.NewBuffer(data))
	CheckError(t, err)

	req.Header.Set(SignatureHeader, string(signature))
	req.Header.Set(PublicKeyHeader, encodedPubKey)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.AddWorker)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	workers, err := s.RetrieveWorkersByQueueID(1)
	CheckError(t, err)
	contains := false
	for _, w := range workers {
		if w.Equals(&worker) {
			contains = true
			break
		}
	}

	if !contains {
		t.Fatal("The worker was not found in the storage")
	}

	var _token map[string]string
	err = json.NewDecoder(rr.Body).Decode(&_token)
	CheckError(t, err)

	returnedToken, ok := _token["arrebol-worker-token"]

	if !ok {
		t.Error("Arrebol token has not been returned")
	}

	if !token.Token(returnedToken).IsValid() {
		t.Error("The token was invalid")
	}
}

func CheckError(t *testing.T, error error) {
	if error != nil {
		t.Fatal(error)
	}
}