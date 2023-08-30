package usecases

import (
	"github.com/devopshub-tech/devopshub/internal/application/dtos"
	"github.com/devopshub-tech/devopshub/internal/application/mappers"
	"github.com/devopshub-tech/devopshub/internal/application/services"
	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
)

type JobUsecase struct {
	jobService *services.JobService
	mapper     *mappers.JobMapper
	logger     domain.ILogger
}

func NewJobUsecase(jobService *services.JobService, mapper *mappers.JobMapper) *JobUsecase {
	return &JobUsecase{
		jobService: jobService,
		mapper:     mapper,
		logger:     logging.NewLogger(),
	}
}

func (u *JobUsecase) CreateAndEnqueueJob(createJobDTO *dtos.CreateJobDTO) (*dtos.JobDTO, error) {
	pluginInfo := domain.PluginInfo{
		Name:    createJobDTO.PluginName,
		Version: createJobDTO.PluginVersion,
	}
	newJob := domain.NewJob(pluginInfo, createJobDTO.Inputs)

	job, err := u.jobService.JobCreationAndQueue(newJob)
	if err != nil {
		u.logger.Errorf("Failed to submit job: %v", err)
		return nil, err
	}

	return u.mapper.ToResponse(job), nil
}

func (u *JobUsecase) GetJobById(id string) (*domain.Job, error) {
	job, err := u.jobService.FindJobById(id)
	if err != nil {
		u.logger.Errorf("Failed to get job by ID %s: %v", id, err)
		return nil, err
	}

	if job == nil {
		return nil, nil
	}

	return job, nil
}
