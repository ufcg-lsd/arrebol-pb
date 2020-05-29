package token

import (
	"fmt"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"strings"
)

type Token interface {
	String() string
}

type SimpleToken struct {
	WorkerId string `json:"WorkerId"`
	QueueId string  `json:"QueueId"`
	VCPU	float32 `json:"vCPU"`
	RAM     uint32  `json:"RAM"`
}

func (t *SimpleToken) String() string {
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

type Generator interface {
	NewToken(worker *worker.Worker) (Token, error)
}

type SimpleGenerator struct {}

func NewSimpleGenerator() Generator {
	return &SimpleGenerator{}
}

func (g *SimpleGenerator) NewToken(worker *worker.Worker) (Token, error) {
	return &SimpleToken{
		WorkerId: worker.ID,
		QueueId:  worker.QueueId,
		VCPU:     worker.VCPU,
		RAM:      worker.RAM,
	}, nil
}
