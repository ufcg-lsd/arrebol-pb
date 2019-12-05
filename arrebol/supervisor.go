package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
	"os"
	"strconv"
	"sync"
)

type Supervisor struct {
	queue            *storage.Queue
	availableWorkers []*Worker
	pendingTasks     chan *storage.Task
	pendingPlans 	 chan *AllocationPlan
	scheduler        *Scheduler
	mux              sync.Mutex
}

func NewSupervisor(queue *storage.Queue) *Supervisor {
	return &Supervisor{
		queue:        queue,
		pendingTasks: make(chan *storage.Task),
		pendingPlans: make(chan *AllocationPlan),
		scheduler:    NewScheduler(Fifo),
	}
}

// should be specific by node
func (s *Supervisor) HireRawWorkers(driver Driver) {
	switch driver {
	case Raw:
		log.Println("just support system level execution with static pool of workers")
		pool, _ := strconv.Atoi(os.Getenv("STATIC_WORKER_POOL"))

		for i := 0; i < pool; i++ {
			s.availableWorkers = append(s.availableWorkers, NewWorker(Raw))
		}

	case Docker:
		log.Println("not supported yet")
	default:
		log.Println("no worker type")
	}
}

// Starts the supervisor protocol with a static default scheduler
func (s *Supervisor) Start() {
	log.Printf("Supervisor of queue < %d > started\n", s.queue.ID)
	// only support raw workers, for now, meaning jobs sent to this supervisor will run
	// uninsulated and on the Unix-type host operating system
	s.HireRawWorkers(Raw)

	s.inferPlans()

	s.pokeScheduler()
}

func (s *Supervisor) Collect(job *storage.Job) {
	log.Printf("Collecting tasks of the job %d", job.ID)
	tasks := &job.Tasks
	for _, task := range *tasks {
		s.pendingTasks <- task
	}
}

type AllocationPlan struct {
	task *storage.Task
	worker *Worker
}

func (s *Supervisor) pokeScheduler() {
	for plan := range s.pendingPlans {
		s.scheduler.Schedule(plan)
	}
}

// Seeding to the channel of plans.
// Listening to the channel of pending tasks.
// Ever that a new task exists this method will be called
// generating a new resource allocation plan to execute the task
func (s *Supervisor) inferPlans() {
	for task := range s.pendingTasks {
		log.Println("new pending task")

		plan := s.inferPlanForTask(task)

		if plan != nil {
			s.pendingPlans <- plan
		}
	}
}

func (s *Supervisor) inferPlanForTask(task *storage.Task) *AllocationPlan {
	s.mux.Lock()
	var w *Worker
	for _, worker := range s.availableWorkers {
		if worker.MatchAny(task) {
			w = worker
		}
	}
	defer s.mux.Unlock()
	if w != nil {
		 return	&AllocationPlan{
			task: task,
			worker: w,
		}
	} else {
		return nil
	}
}
