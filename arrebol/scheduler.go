package arrebol

// no preemptive
type Scheduler struct {
	policy Policy
}

type Policy uint

const (
	Fifo Policy = iota
)

func (p Policy) String() string {
	return [...]string{"Fifo"}[p]
}

func (p Policy) Schedule(plan *AllocationPlan) {
	switch p {
	case Fifo:

	}
}

func NewScheduler(policy Policy) *Scheduler {
	return &Scheduler{
		policy: policy,
	}
}

func (s *Scheduler) Schedule(plan *AllocationPlan) {
	s.policy.Schedule(plan)
}

