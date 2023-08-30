// cmd/gateway/main.go
package main

import (
	"flag"
	"fmt"

	apphttp "github.com/devopshub-tech/devopshub/internal/application/http"
	"github.com/devopshub-tech/devopshub/internal/application/mappers"
	"github.com/devopshub-tech/devopshub/internal/application/services"
	"github.com/devopshub-tech/devopshub/internal/application/usecases"
	infconfig "github.com/devopshub-tech/devopshub/internal/infrastructure/config"
	infhttp "github.com/devopshub-tech/devopshub/internal/infrastructure/http"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/mongodb"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/queue"
)

var (
	logger = logging.NewLogger()
	config = infconfig.NewConfig()
	host   = flag.String("host", config.GetApiHost(), fmt.Sprintf("api host, default: %s", config.GetApiHost()))
	port   = flag.Int("port", config.GetApiPort(), fmt.Sprintf("api host, default: %d", config.GetApiPort()))
)

func main() {
	flag.Parse()

	logger.Info("Connecting to database...")
	db, err := mongodb.NewMongoDB(config.GetDbMongoUri(), config.GetDbName())
	if err != nil {
		logger.Fatalf("Error connecting database: %v", err)
	}

	jobMapper := mappers.NewJobMapper()

	jobRepo, err := mongodb.NewJobRepository(db.GetDatabase(), jobMapper)
	if err != nil {
		logger.Fatalf("Error creating job repository: %v", err)
	}

	queue, err := queue.NewRabbitMQ(config.GetQueueUri())
	if err != nil {
		logger.Fatalf("Error connecting service queue: %v", err)
	}

	jobService := services.NewJobService(jobRepo, queue)

	jobUsecase := usecases.NewJobUsecase(jobService, jobMapper)

	router := infhttp.NewGinRouter(config.GetEnvMode())
	router.AddReadinessCheck(queue.CheckQueueHealth)
	router.AddReadinessCheck(db.CheckMongoDBHealth)
	router.SetupHealthRoutes()

	jobRouter := apphttp.NewJobRouter(jobUsecase)
	jobRouter.Setup(router.Engine())

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)
	server := infhttp.NewHTTPServer()
	server.Start(httpAddr, "gateway", router.Engine())
}
