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
	PendingTasks map[uint][]*storage.Task
	ReportInterval int64
	QM *QueuesManager
	S *storage.Storage
}

func NewJobsHandler(s *storage.Storage, q *QueuesManager) JobsHandler{
	reportInterval := os.Getenv(ReportIntervalKey)
	parsedInterval, err := strconv.ParseInt(reportInterval, 10, 64)

	if err != nil {
		log.Fatal("Invalid " + ReportIntervalKey + "env variable")
	}

	return JobsHandler{ReportInterval: parsedInterval, S: s, QM: q}
}

func (j *JobsHandler) Start() {
	//make them all never dieing
	go j.extractPendingTasks()
	go j.checkPendingTasks()
	go j.checkNeverEndingTasks()
	go j.jobsStateChanger()
}

func (j *JobsHandler) extractPendingTasks() {
	for {
		queues := j.QM.GetQueues()

		for _, queue := range queues {
			tasks := j.S.RetrieveTasksByState(queue.ID, storage.TaskPending)
			for _, task := range tasks {
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
					j.PendingTasks[queue.ID] = append(j.PendingTasks[queue.ID], task)
				}
			}
		}

		time.Sleep(SLEEPING_TIME*time.Second)
	}
}

func(j *JobsHandler) GetPendingTasks(queueId uint) []*storage.Task{
	return j.PendingTasks[queueId]
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
		queues := j.QM.GetQueues()

		for _, queue := range queues {
			tasks := j.S.RetrieveTasksByState(queue.ID, storage.TaskRunning)
			for _, task := range tasks {
				var expectedDelay int64 = 30
				if task.UpdatedAt.Unix() + task.ReportInterval > time.Now().Unix() + expectedDelay {
					task.State = storage.TaskPending
					task.Progress = 0
					j.S.SaveTask(task)
					j.PendingTasks[queue.ID] = append(j.PendingTasks[queue.ID], task)
				}
			}
		}

		time.Sleep(SLEEPING_TIME*time.Second)
	}
}

func (j *JobsHandler) checkPendingTasks() {
	for {
		for queueId, tasks := range j.PendingTasks {
			for i, task := range tasks {
				task, _ = j.S.RetrieveTask(task.ID)
				if task.State != storage.TaskPending {
					j.PendingTasks[queueId] = append(j.PendingTasks[queueId][:i], j.PendingTasks[queueId][i+1:]...)
				}
			}
		}

		time.Sleep(SLEEPING_TIME*time.Second)
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

		time.Sleep(SLEEPING_TIME*time.Second)
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
		if !contains(t.State, states) {
			return false
		}
	}
	return true
}

func contains(e storage.TaskState, arr []storage.TaskState) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}