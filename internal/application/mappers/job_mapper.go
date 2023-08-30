// internal/application/mappers/job_mapper.go
package mappers

import (
	"github.com/devopshub-tech/devopshub/internal/application/dtos"
	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/mongodb/models"
)

type JobMapper struct{}

func NewJobMapper() *JobMapper {
	return &JobMapper{}
}

// Mapping from Job model of persistence to Job of domain
func (m *JobMapper) ToDomain(jobModel *models.JobModel) *domain.Job {
	if jobModel == nil {
		return nil
	}

	return &domain.Job{
		Id:        jobModel.Id,
		Plugin:    domain.PluginInfo(jobModel.Plugin),
		Status:    jobModel.Status,
		Log:       jobModel.Log,
		Inputs:    jobModel.Inputs,
		CreatedAt: jobModel.CreatedAt,
		UpdatedAt: jobModel.UpdatedAt,
	}
}

// Map the Job of the domain to the Job model of the persistence
func (m *JobMapper) ToPersistence(job *domain.Job) *models.JobModel {
	if job == nil {
		return nil
	}

	return &models.JobModel{
		Id:        job.Id,
		Plugin:    models.PluginInfoModel(job.Plugin),
		Status:    job.Status,
		Log:       job.Log,
		Inputs:    job.Inputs,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

// Map from Domain Job to Response Job
func (m *JobMapper) ToResponse(job *domain.Job) *dtos.JobDTO {
	return &dtos.JobDTO{
		Id: job.Id,
		Plugin: dtos.PluginInfoDTO{
			Name:    job.Plugin.Name,
			Version: job.Plugin.Version,
		},
		Status:    job.Status,
		Inputs:    job.Inputs,
		Log:       job.Log,
		CreatedAt: job.CreatedAt,
	}
}
