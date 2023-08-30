// internal/application/dtos/job_dtos.go
package dtos

import (
	"time"

	"github.com/devopshub-tech/devopshub/internal/domain"
)

type CreateJobDTO struct {
	PluginName    string            `json:"pluginName" form:"pluginName" validate:"required"`
	PluginVersion string            `json:"pluginVersion" form:"pluginVersion" validate:"required"`
	Inputs        map[string]string `json:"inputs"`
}

type JobDTO struct {
	Id        string            `json:"id"`
	Plugin    PluginInfoDTO     `json:"plugin"`
	Status    domain.JobStatus  `json:"status"`
	Inputs    map[string]string `json:"inputs"`
	Log       string            `json:"log"`
	CreatedAt time.Time         `json:"createdAt"`
}

type PluginInfoDTO struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
