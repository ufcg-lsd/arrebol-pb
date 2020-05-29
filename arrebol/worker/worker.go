package worker

type Worker struct {
	ID      string  `json:"ID"`
	VCPU    float32 `json:"vCPU"`
	RAM     uint32  `json:"RAM"` //Megabytes
	QueueId string  `json:"QueueId, omitempty"`
}
