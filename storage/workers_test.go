package storage

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"testing"
	"time"
)

func TestHasUUID(t *testing.T) {
	s := OpenDriver()
	s.CreateTables()

	t.Run("assert that the insertion works", func(t *testing.T) {
		if !s.driver.HasTable(&worker.Worker{}) {
			t.Errorf("expected has a table but nothing was found")

			expected := uuid.NewV4()

			w := worker.Worker{
				Base: worker.Base{
					ID:        expected,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: nil,
				},
				VCPU:    0,
				RAM:     0,
				QueueID: 0,
			}

			has, err := s.SaveWorker(w)

			if err != nil && has != uuid.Nil {
				t.Errorf("nil found")
			}

			if has != expected && has != uuid.Nil {
				t.Errorf("expected %s but has %s", expected.String(), has.String())
			}
		}
	})

	err := s.driver.Close()

	if err != nil {
		t.Errorf("error %T", err)
	}
}
