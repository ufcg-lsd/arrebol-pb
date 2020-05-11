package arrebol

import (
	"fmt"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/driver"
	"github.com/ufcg-lsd/arrebol-pb/docker"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// no preemptive
type Scheduler struct {
	workers      []*Worker
	pendingTasks chan *storage.Task
	pendingPlans chan *AllocationPlan
	policy       Policy
	mutex        sync.Mutex
}

type Policy uint

const TaskRetryTimeInterval = 10 * time.Second

const (
	Fifo Policy = iota
)

func (p Policy) String() string {
	return [...]string{"Fifo"}[p]
}

func (p Policy) schedule(plans chan *AllocationPlan) {
	switch p {
	case Fifo:
		for plan := range plans {
			go plan.execute()
		}
	default:
		log.Println("Just support fifo")
	}
}

func NewScheduler(policy Policy) *Scheduler {
	return &Scheduler{
		policy:       policy,
		workers:      make([]*Worker, 0),
		pendingTasks: make(chan *storage.Task),
		pendingPlans: make(chan *AllocationPlan),
	}
}

func (s *Scheduler) Start() {
	// only support raw workers, for now, meaning jobs sent to the supervisor of this scheduler will run
	// uninsulated and on the Unix-type host operating system
	s.HireWorkers()
	go s.inferPlans()
	s.Schedule()
}

func (s *Scheduler) Schedule() {
	s.policy.schedule(s.pendingPlans)
}

// should be specific by node
func (s *Scheduler) HireWorkers() {
	pool, _ := strconv.Atoi(os.Getenv("WORKERS_AMOUNT"))
	if os.Getenv("DRIVER") == "docker" {
		address := os.Getenv("WORKER_ADDRESS")
		cli := docker.NewDockerClient(address)
		for i := 0; i < pool; i++ {
			_driver := driver.DockerDriver{
				Id:  fmt.Sprintf("docker-worker-%d", i),
				Cli: *cli,
			}
			s.workers = append(s.workers, NewWorker(&_driver))
		}
	} else {
		_driver := driver.RawDriver{}
		for i := 0; i < pool; i++ {
			s.workers = append(s.workers, NewWorker(&_driver))
		}
	}
	//log.Println("just support system level execution with static pool of workers")
}

func (s *Scheduler) AddTask(task *storage.Task) {
	s.pendingTasks <- task
}

type AllocationPlan struct {
	task *storage.Task
	worker *Worker
}

func (a *AllocationPlan) execute() {
	a.worker.Execute(a.task)
}

// Seeding to the channel of plans.
// Listening to the channel of pending tasks.
// Ever that a new task exists this method will be called
// generating a new resource allocation plan to execute the task
func (s *Scheduler) inferPlans() {
	for {
		task := <- s.pendingTasks
		log.Printf("Planning to run task [%d]", task.ID)

		plan := s.inferPlanForTask(task)

		if plan != nil {
			s.pendingPlans <- plan // a channel is used here because only fifo's policy is supported
		} else {
			go func() {
				time.Sleep(TaskRetryTimeInterval)
				s.pendingTasks <- task
				log.Printf("Retring the task [%d]", task.ID)
			}()
		}
	}

}

func (s *Scheduler) inferPlanForTask(task *storage.Task) *AllocationPlan {
	s.mutex.Lock()
	log.Printf("Searching worker for task [%d]", task.ID)
	for _, worker := range s.workers {
		if worker.MatchAny(task) {
			log.Printf("The task [%d] matched with the worker [%s]", task.ID, worker.id)
			return s.makePlan(worker, task)
		}
	}
	defer s.mutex.Unlock()
	return nil
}

func (s *Scheduler) makePlan(w *Worker, t *storage.Task) *AllocationPlan {
	w.state = Busy
	// TODO Change task state to pending or queued
	return &AllocationPlan{
		task: t,
		worker: w,
	}
}

