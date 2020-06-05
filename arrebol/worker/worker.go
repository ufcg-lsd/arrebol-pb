package worker

import "github.com/jinzhu/gorm"

type Worker struct {
	gorm.Model
	WorkerID string  `json:"Id"`
	VCPU     float32 `json:"Vcpu"`
	RAM      uint32  `json:"Ram"` //Megabytes
	QueueID  uint  `json:"QueueID, omitempty"`
}
