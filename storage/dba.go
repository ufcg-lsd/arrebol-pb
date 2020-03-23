package storage

import "errors"

func (t *Task) GetRawCommands() []string {
	var raws []string
	for _, c := range t.Commands {
		raws = append(raws, c.RawCommand)
	}
	return raws
}

func (t *Task) GetConfig(key string) (string, error) {
	for _, conf := range t.Config {
		if conf.Key == key {
			return conf.Value, nil
		}
	}
	return "", errors.New("Config [" + key + "] not found")
}

func (q Queue) QueueHasJob(jobId uint) bool {
	for _, job := range q.Jobs {
		if job.ID == jobId {
			return true
		}
	}
	return false
}