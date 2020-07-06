package worker

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Worker struct {
	Base
	VCPU    float32 `json:"Vcpu"`
	RAM     uint32  `json:"Ram"` //Megabytes
	QueueID uint    `json:"QueueID, omitempty"`
}

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (w *Worker) Equals(o *Worker) bool {
	if o != nil && w.ID == o.ID {
		return true
	}
	return false
}
