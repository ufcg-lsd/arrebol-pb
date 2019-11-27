package api

import (
	"github.com/emanueljoivo/arrebol/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func responseFromQueue(queue *storage.Queue) *QueueResponse {
	var pendingTasks uint
	var runningTasks uint
	jobs := queue.Jobs

	for i := 0; i < len(jobs); i++ {
		tasks := jobs[i].Tasks
		for j := 0; j< len(tasks); j++ {
			if tasks[j].State == storage.Pending {
				pendingTasks++
			}
			if tasks[j].State == storage.Running {
				runningTasks++
			}
		}
	}

	return &QueueResponse{
		ID: queue.ID.Hex(),
		Name: queue.Name,
		PendingTasks: pendingTasks,
		RunningTasks: runningTasks,
		Nodes: uint(len(queue.Nodes)),
	}
}

func jobFromSpec(jobSpec JobSpec, id primitive.ObjectID) *storage.Job {
	var tasks []storage.Task

	taskSpecs := jobSpec.Tasks

	for i := 0; i< len(taskSpecs); i++ {
		tasks = append(tasks, storage.Task{
			ID:       taskSpecs[i].ID,
			Config:   taskSpecs[i].Config,
			Commands: taskSpecs[i].Commands,
			State:    storage.Pending,
		})
	}

	return &storage.Job{
		ID:    id,
		Label: jobSpec.Label,
		State: storage.Pending,
		Tasks: tasks,
	}
}
