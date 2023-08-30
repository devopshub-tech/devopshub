// internal/application/http/job_router.go
package http

import (
	"github.com/devopshub-tech/devopshub/internal/application/usecases"
	"github.com/gin-gonic/gin"
)

type JobRouter struct {
	handlers *Handlers
}

func NewJobRouter(jobUsecase *usecases.JobUsecase) *JobRouter {
	handlers := NewHandlers(jobUsecase)
	return &JobRouter{handlers: handlers}
}

func (r *JobRouter) Setup(router *gin.Engine) {
	router.GET("/job/:id", r.handlers.GetJobDetails)
	router.POST("/job", r.handlers.CreateJob)
}
