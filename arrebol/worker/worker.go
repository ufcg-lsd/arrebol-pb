package worker

type Worker struct {
	ID      string  `json:"Id"`
	VCPU    float32 `json:"Vcpu"`
	RAM     uint32  `json:"Ram"` //Megabytes
	QueueID uint    `json:"QueueID, omitempty"`
}

func (w *Worker) Equals(o *Worker) bool {
	if o != nil && w.ID == o.ID {
		return true
	}
	return false
}