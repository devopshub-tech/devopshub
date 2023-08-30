// internal/application/queue/worker_router.go
package queue

import (
	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	infrqueue "github.com/devopshub-tech/devopshub/internal/infrastructure/queue"
)

type WorkerRouter struct {
	routers domain.QueueRouter
	logger  domain.ILogger
}

func NewWorkerRouter() *WorkerRouter {
	return &WorkerRouter{
		routers: make(domain.QueueRouter),
		logger:  logging.NewLogger(),
	}
}

func (r *WorkerRouter) Setup() {
	r.routers[infrqueue.JobQueueName] = func(message []byte) error {
		// Simulation process job
		job := string(message)
		r.logger.Infof("Processing job: %s\n", job)
		return nil
	}
}

func (r *WorkerRouter) Routers() domain.QueueRouter {
	return r.routers
}
