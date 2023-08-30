package domain

import (
	"time"

	"github.com/devopshub-tech/devopshub/pkg/utils"
)

// JobStatus represents the possible statuses of a job.
type JobStatus string

const (
	JobStatusCreated   JobStatus = "created"
	JobStatusEnqueued  JobStatus = "enqueued"
	JobStatusExecuting JobStatus = "executing"
	JobStatusFinished  JobStatus = "finished"
	JobStatusFailed    JobStatus = "failed"
)

// Job represents a job requested from the application.
type Job struct {
	Id        string
	Plugin    PluginInfo
	Status    JobStatus
	Log       string
	Inputs    map[string]string // Variables for plugin execution
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewJob creates a new instance of Job.
func NewJob(plugin PluginInfo, inputs map[string]string) *Job {
	return &Job{
		Id:        utils.GenerateUUID(),
		Plugin:    plugin,
		Status:    JobStatusCreated,
		Inputs:    inputs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetStatus sets the status of a job.
func (j *Job) SetStatus(status JobStatus) {
	j.Status = status
}

// AddLog adds information to a job log.
func (j *Job) AgregarLog(log string) {
	j.Log += log
}

// JobRepository defines the interface to interact with the job repository.
type IJobRepository interface {
	Create(job *Job) (interface{}, error)
	Update(job *Job) error
	Find(filter interface{}) ([]*Job, error)
	FindById(id string) (*Job, error)
}
