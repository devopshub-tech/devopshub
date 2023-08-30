// cmd/worker/main.go
package main

import (
	"flag"
	"fmt"

	appqueue "github.com/devopshub-tech/devopshub/internal/application/queue"
	infrconfig "github.com/devopshub-tech/devopshub/internal/infrastructure/config"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	infrqueue "github.com/devopshub-tech/devopshub/internal/infrastructure/queue"
)

var (
	logger  = logging.NewLogger()
	config  = infrconfig.NewConfig()
	rmqHost = flag.String("host", config.GetQueueUri(), fmt.Sprintf("queue host, default: %s", config.GetQueueUri()))
)

func main() {
	flag.Parse()

	queue, err := infrqueue.NewRabbitMQ(config.GetQueueUri())
	if err != nil {
		logger.Fatalf("Error connecting service queue: %v", err)
	}

	router := appqueue.NewWorkerRouter()
	router.Setup()

	server := infrqueue.NewQueueServer(queue)
	server.Start("worker", router.Routers())
}
