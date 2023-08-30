// internal/infrastructure/queue/server.go
package queue

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
)

type QueueServer struct {
	engine domain.IQueue
	logger domain.ILogger
}

func NewQueueServer(engine domain.IQueue) *QueueServer {
	return &QueueServer{
		engine: engine,
		logger: logging.NewLogger(),
	}
}

func (qs *QueueServer) Start(serviceName string, router domain.QueueRouter) {
	if len(router) == 0 {
		qs.logger.Infof("No queues configured for '%s'. Exiting queue server.", serviceName)
		return
	}

	qs.startQueueListeners(serviceName, router)
	qs.startHealthChecks()

	qs.waitForShutdown(serviceName)
}

func (qs *QueueServer) startQueueListeners(serviceName string, router domain.QueueRouter) {
	for name, handler := range router {
		go qs.engine.Consume(name, fmt.Sprint(serviceName, "-", name), handler)
		qs.logger.Infof("Queue '%s' started listening.", name)
	}
	qs.logger.Infof("Queue server '%s' started.", serviceName)
}

func (qs *QueueServer) startHealthChecks() {
	qs.manageHealthFile("liveness.txt", "create", "UP")
	qs.manageHealthFile("readiness.txt", "create", "NOT_READY")
	qs.SetupReadinessCheck(10 * time.Second)
}

func (qs *QueueServer) waitForShutdown(serviceName string) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	qs.logger.Infof("Shutting down queue server '%s'...", serviceName)

	qs.manageHealthFile("liveness.txt", "remove", "DOWN")
	qs.manageHealthFile("readiness.txt", "remove", "NOT_READY")

	if err := qs.engine.Close(); err != nil {
		qs.logger.Errorf("Error closing queue engine: %v", err)
	}

	qs.logger.Infof("Queue server '%s' exited.", serviceName)
}

func (qs *QueueServer) manageHealthFile(healthFileName, action, message string) {
	livenessFilePath := filepath.Join(os.TempDir(), healthFileName)

	switch action {
	case "create":
		createHealthFile(livenessFilePath, message, qs.logger)
	case "remove":
		removeHealthFile(livenessFilePath, qs.logger)
	}
}

func createHealthFile(filePath, message string, logger domain.ILogger) {
	livenessFile, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("Error creating %s file: %v", filePath, err)
		return
	}
	defer livenessFile.Close()

	livenessFile.WriteString(message)
	logger.Debugf("Health %s file created.", filePath)
}

func removeHealthFile(filePath string, logger domain.ILogger) {
	err := os.Remove(filePath)
	if err != nil {
		logger.Errorf("Error removing %s file: %v", filePath, err)
		return
	}

	logger.Debugf("Health %s file removed.", filePath)
}

func (qs *QueueServer) checkQueueHealthPeriodically(interval time.Duration, resultChan chan<- bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if qs.engine.CheckQueueHealth() {
				resultChan <- true
			} else {
				resultChan <- false
			}
		}
	}
}

func (qs *QueueServer) SetupReadinessCheck(interval time.Duration) {
	readinessChan := make(chan bool)

	go qs.checkQueueHealthPeriodically(interval, readinessChan)

	go func() {
		for {
			select {
			case isHealthy := <-readinessChan:
				if isHealthy {
					qs.manageHealthFile("readiness.txt", "create", "READY")
				} else {
					qs.manageHealthFile("readiness.txt", "remove", "NOT_READY")
				}
			}
		}
	}()
}
