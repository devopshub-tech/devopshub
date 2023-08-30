// internal/infrastructure/queue/rabbitmq.go
package queue

import (
	"time"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/config"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  domain.IConfig
	logger  domain.ILogger
}

// NewRabbitMQ creates a new instance of RabbitMQ.
func NewRabbitMQ(connectionString string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		config:  config.NewConfig(),
		logger:  logging.NewLogger(),
	}, nil
}

// Publish publishes a message to a queue.
func (q *RabbitMQ) Publish(queueName string, message []byte) error {
	err := q.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		q.logger.Debugf("Error publishing message: %v", err)
		return err
	}

	return nil
}

// Enqueue enqueues a message to a queue (alias for Publish).
func (q *RabbitMQ) Enqueue(queueName string, message []byte) error {
	return q.Publish(queueName, message)
}

// Consume starts consuming messages from a queue with the provided handleq.
func (q *RabbitMQ) Consume(queueName string, consumerName string, handler domain.QueueHandler) error {
	msgs, err := q.channel.Consume(
		queueName,
		consumerName,
		true, // Auto-acknowledgement
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		q.logger.Debugf("Error consuming messages: %v", err)
		return err
	}

	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				q.logger.Debugf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}

// Close closes the RabbitMQ connection and channel.
func (q *RabbitMQ) Close() error {
	err := q.channel.Close()
	if err != nil {
		q.logger.Debugf("Error closing channel: %v", err)
		return err
	}

	err = q.conn.Close()
	if err != nil {
		q.logger.Debugf("Error closing connection: %v", err)
		return err
	}

	q.logger.Info("Connection closed with queues successful.")
	return nil
}

func (q *RabbitMQ) CheckQueueHealth() bool {
	queueName := "healthy"
	consumerTag := "healthy-consumer"

	_, err := q.channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		q.logger.Errorf("Error declaring queue:", err)
		return false
	}

	message := "Hello, RabbitMQ!"
	err = q.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	if err != nil {
		q.logger.Errorf("Error posting healthy message:", err)
		return false
	}

	msgs, err := q.channel.Consume(queueName, consumerTag, true, false, false, false, nil)
	if err != nil {
		q.logger.Errorf("Error consuming healthy message:", err)
		return false
	}

	select {
	case <-msgs:
		q.channel.Cancel(consumerTag, false)
		return true
	case <-time.After(time.Duration(q.config.GetQueueHealthTimeout()) * time.Second):
		q.logger.Error("Timed out waiting to consume the message.")
		q.channel.Cancel(consumerTag, false)
		return false
	}
}
