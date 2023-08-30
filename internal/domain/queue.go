// internal/domain/queue.go
package domain

type QueueHandler func(message []byte) error

type QueueRouter map[string]QueueHandler

type IQueue interface {
	Publish(queueName string, message []byte) error
	Enqueue(queueName string, message []byte) error
	Consume(queueName string, consumerName string, handler QueueHandler) error
	CheckQueueHealth() bool
	Close() error
}
