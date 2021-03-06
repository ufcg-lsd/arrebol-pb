package service

import (
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"os"
	"strconv"
	"time"
)

const ReportIntervalKey = "REPORT_INTERVAL"

// in seconds
const SLEEPING_TIME = 5

type JobsHandler struct {
	PendingTasks   map[uint][]*storage.Task
	ReportInterval int64
	S              *storage.Storage
	opChan chan Operation
}

func NewJobsHandler(s *storage.Storage) *JobsHandler {
	reportInterval := os.Getenv(ReportIntervalKey)
	parsedInterval, err := strconv.ParseInt(reportInterval, 10, 64)

	if err != nil {
		log.Fatal("Invalid " + ReportIntervalKey + "env variable")
	}

	return &JobsHandler{ReportInterval: parsedInterval, S: s,
		opChan: make(chan Operation, 1000), PendingTasks: map[uint][]*storage.Task{}}
}

type OperationType string

const (
	OperationAdd OperationType = "ADD"
	OperationRemove OperationType = "REMOVE"
)


func (j *JobsHandler) Start() {
	log.Println("Starting jobs handler")
	//The aggregateTasks func receives messages through the opChan
	//and modify the pendingTasks according to the message content.
	//This avoid some concurrency issues that would happen if the
	//other goroutines modified the pendingTasks by their own.
	go j.aggregateTasks()
	go j.extractPendingTasks()
	go j.checkNeverEndingTasks()
	go j.jobsStateChanger()
}

type OperationValue struct {
	qid uint
	t []*storage.Task
}

type Operation struct {
	Type OperationType
	v OperationValue
}

func (j *JobsHandler) aggregateTasks() {
	log.Println("Starting tasks agregator")
	for {
		select {
		case op := <-j.opChan:
			log.Println("Received operation: " + op.Type)
			log.Println("Tasks size: " + strconv.Itoa(len(op.v.t)) + "from queue: " + string(op.v.qid))
			switch op.Type {
			case OperationAdd:
				j.PendingTasks[op.v.qid] = append(j.PendingTasks[op.v.qid], op.v.t...)
			case OperationRemove:
				for _, t := range op.v.t {
					j.PendingTasks[op.v.qid] = removeTask(t, j.PendingTasks[op.v.qid])
				}
			}
		}
	}
}

func removeTask(task *storage.Task, tasks []*storage.Task) []*storage.Task{
	returningTasks := []*storage.Task{}
	for i, t := range tasks {
		if t == task {
			returningTasks = append(tasks[:i], tasks[i+1:]...)
		}
	}
	return returningTasks
}

func getPendingTasksFromQueue(q *storage.Queue) []*storage.Task{
	tasks := []*storage.Task{}
	log.Println("Queue jobs:")
	log.Println(q.Jobs)
	for _, job := range q.Jobs {
		if job.State == storage.JobFailed || job.State == storage.JobFinished {
			continue
		}
		log.Println(job.Tasks)

		for _, task := range job.Tasks {
			if task.State == storage.TaskPending {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

func qAsStr(queues []*storage.Queue) string {
	str := ""
	for _, q := range queues {
		str += " " + q.Name
	}
	return str
}

func (j *JobsHandler) extractPendingTasks() {
	for {
		queues := loadQueues(j.S)

		log.Println("Extracting pending tasks from queues " + qAsStr(queues))

		for _, queue := range queues {
			tasks := getPendingTasksFromQueue(queue)
			log.Println("pending tasks:")
			log.Println(tasks)
			tasksToAdd := []*storage.Task{}
			for _, task := range tasks {
				log.Println(task.ID)
				// prevent duplicates
				found := false
				for _, t := range j.PendingTasks[queue.ID] {
					if t == task {
						found = true
						break
					}
				}

				if !found {
					task.ReportInterval = j.ReportInterval
					j.S.SaveTask(task)
					tasksToAdd = append(tasksToAdd, task)
				}
			}
			log.Println("tasks to add:")
			log.Println(tasksToAdd)
			if len(tasksToAdd) > 0 {
				j.opChan <- Operation{
					Type: OperationAdd,
					v: OperationValue{queue.ID, tasksToAdd},
				}
			}
		}

		time.Sleep(SLEEPING_TIME * time.Second)
	}
}

func (j *JobsHandler) GetPendingTasks(queueId uint) []*storage.Task {
	tasks, ok := j.PendingTasks[queueId]

	if !ok {
		return []*storage.Task{}
	}

	j.opChan <- Operation{
		Type: OperationRemove,
		v:    OperationValue{queueId, tasks},
	}

	log.Println("Send pending tasks:")
	log.Println(tasks)

	return tasks
}

func (j *JobsHandler) HandleReport(task *storage.Task) error {
	retrievedTask, err := j.S.RetrieveTask(task.ID)

	if err != nil {
		return err
	}

	retrievedTask.Progress = task.Progress
	retrievedTask.State = task.State
	j.S.SaveTask(task)

	return nil
}

func (j *JobsHandler) checkNeverEndingTasks() {
	for {
		queues := loadQueues(j.S)

		for _, queue := range queues {
			tasks := j.S.RetrieveTasksFromQueueByState(queue.ID, storage.TaskRunning)
			for _, task := range tasks {
				var expectedDelay int64 = 30
				if task.UpdatedAt.Unix()+task.ReportInterval > time.Now().Unix()+expectedDelay {
					task.State = storage.TaskPending
					task.Progress = 0
					j.S.SaveTask(task)
				}
			}
		}

		time.Sleep(SLEEPING_TIME * time.Second)
	}
}

func (j *JobsHandler) jobsStateChanger() {
	for {
		jobs, err := j.S.RetrieveJobs()

		if err != nil {
			log.Fatal("Error on retrieving jobs. " + err.Error())
		}

		for _, job := range jobs {
			jobState := inferJobState(job)
			if jobState != job.State {
				job.State = jobState
				j.S.SaveJob(job)
			}
		}

		time.Sleep(SLEEPING_TIME * time.Second)
	}
}

func inferJobState(job *storage.Job) storage.JobState {
	var jobState storage.JobState
	if isAll([]storage.TaskState{storage.TaskFailed}, job.Tasks) {
		jobState = storage.JobFailed
	} else if isAll([]storage.TaskState{storage.TaskFailed, storage.TaskFinished}, job.Tasks) {
		jobState = storage.JobFinished
	} else if isAll([]storage.TaskState{storage.TaskPending}, job.Tasks) {
		jobState = storage.JobQueued
	} else {
		jobState = storage.JobRunning
	}
	return jobState
}

func isAll(states []storage.TaskState, tasks []*storage.Task) bool {
	for _, t := range tasks {
		if !_contains(t.State, states) {
			return false
		}
	}
	return true
}

func _contains(e storage.TaskState, arr []storage.TaskState) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}
