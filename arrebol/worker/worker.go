package worker

type Worker struct {
	ID      string  `json:"Id"`
	VCPU    float32 `json:"Vcpu"`
	RAM     uint32  `json:"Ram"` //Megabytes
	QueueID uint    `json:"QueueID, omitempty"`
}
