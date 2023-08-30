package services

import (
	"encoding/json"
	"fmt"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/queue"
)

type JobService struct {
	jobRepository domain.IJobRepository
	queue         domain.IQueue
	logger        domain.ILogger
}

func NewJobService(jobRepo domain.IJobRepository, queue domain.IQueue) *JobService {
	return &JobService{
		jobRepository: jobRepo,
		queue:         queue,
		logger:        logging.NewLogger(),
	}
}

func (s *JobService) JobCreationAndQueue(job *domain.Job) (*domain.Job, error) {
	s.logger.Info("Creating and enqueuing new job...")

	job.Status = domain.JobStatusCreated

	_, err := s.jobRepository.Create(job)
	if err != nil {
		s.logger.Errorf("Failed to create job in repository: %v", err)
		return nil, err
	}

	job.Status = domain.JobStatusEnqueued
	jobJSON, err := json.Marshal(job)
	if err != nil {
		errHandler := s.handleCreateJobError("Failed to serialize job", job, err)
		if errHandler != nil {
			return nil, errHandler
		}
		return nil, err
	}

	err = s.queue.Enqueue(queue.JobQueueName, jobJSON)
	if err != nil {
		errHandler := s.handleCreateJobError("Failed to enqueue job", job, err)
		if errHandler != nil {
			return nil, errHandler
		}
		return nil, err
	}

	err = s.jobRepository.Update(job)
	if err != nil {
		s.logger.Errorf("Failed to update job in repository: %v", err)
		return nil, err
	}

	s.logger.Infof("Job created and enqueued successfully. ID: %s", job.Id)

	return job, nil
}

func (s *JobService) FindJobById(id string) (*domain.Job, error) {
	s.logger.Infof("Fetching job details for ID: %s...", id)

	job, err := s.jobRepository.FindById(id)
	if err != nil {
		s.logger.Errorf("Failed to find job by ID: %s, error: %v", id, err)
		return nil, err
	}

	if job == nil {
		s.logger.Infof("Job not found. ID: %s", id)
		return nil, nil
	}

	s.logger.Infof("Job details fetched successfully. ID: %s", job.Id)

	return job, nil
}

func (s *JobService) handleCreateJobError(message string, job *domain.Job, err error) error {
	s.logger.Errorf("Error during job creation: %v: %v", message, err)

	job.Status = domain.JobStatusFailed
	job.Log = fmt.Sprintf("%v: %v", message, err)
	updateErr := s.jobRepository.Update(job)
	if updateErr != nil {
		s.logger.Errorf("Failed to update job status after error: %v", updateErr)
		return updateErr
	}

	return err
}
