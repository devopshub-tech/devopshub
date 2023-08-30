// internal/application/http/job_handlers.go
package http

import (
	"fmt"
	"net/http"

	"github.com/devopshub-tech/devopshub/internal/application/dtos"
	"github.com/devopshub-tech/devopshub/internal/application/usecases"
	"github.com/devopshub-tech/devopshub/internal/application/utils"
	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	jobUsecase *usecases.JobUsecase
	logger     domain.ILogger
}

func NewHandlers(jobUsecase *usecases.JobUsecase) *Handlers {
	return &Handlers{
		jobUsecase: jobUsecase,
		logger:     logging.NewLogger(),
	}
}

func (h *Handlers) CreateJob(ctx *gin.Context) {
	var createJobDTO *dtos.CreateJobDTO

	if err := ctx.ShouldBindJSON(&createJobDTO); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		RespondError(ctx.Writer, "Invalid input data.", http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(createJobDTO); err != nil {
		h.logger.Errorf("Error validating fields: %v", err)
		RespondError(ctx.Writer, "Error validating fields.", http.StatusBadRequest, utils.FormatErrors(err)...)
		return
	}

	job, err := h.jobUsecase.CreateAndEnqueueJob(createJobDTO)
	if err != nil {
		h.logger.Errorf("Failed to submit job: %v", err)
		RespondError(ctx.Writer, "Failed to submit job", http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(ctx.Writer, job, "Job created successfully", http.StatusOK)
}

func (h *Handlers) GetJobDetails(ctx *gin.Context) {
	id := ctx.Param("id")

	job, err := h.jobUsecase.GetJobById(id)
	if err != nil {
		h.logger.Errorf("Failed to fetch job details for ID %s: %v", id, err)
		RespondError(ctx.Writer, fmt.Sprintf("Failed to fetch job. ID %s", id), http.StatusInternalServerError, err.Error())
		return
	}

	if job == nil {
		h.logger.Errorf("Failed to fetch job details for ID %s: %v", id, err)
		RespondError(ctx.Writer, "Job not found", http.StatusNotFound)
		return
	}

	RespondSuccess(ctx.Writer, job, "Job details retrieved successfully", http.StatusOK)
}
