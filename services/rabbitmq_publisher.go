package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/icpinto/dating-app/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultLifecycleExchange = "user.lifecycle"
	lifecycleRoutingTemplate = "user.%s"
)

// RabbitMQPublisher wraps a RabbitMQ connection/channel for publishing lifecycle events.
type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// NewRabbitMQPublisher establishes a connection to RabbitMQ and declares the lifecycle exchange.
func NewRabbitMQPublisher(url string, exchange string) (*RabbitMQPublisher, error) {
	if exchange == "" {
		exchange = defaultLifecycleExchange
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}
	return &RabbitMQPublisher{conn: conn, channel: ch, exchange: exchange}, nil
}

// PublishLifecycleEvent publishes the lifecycle event payload to the configured exchange.
func (p *RabbitMQPublisher) PublishLifecycleEvent(ctx context.Context, event models.UserLifecycleOutbox) error {
	if p == nil || p.channel == nil {
		return fmt.Errorf("rabbitmq publisher not configured")
	}
	routingKey := fmt.Sprintf(lifecycleRoutingTemplate, string(event.EventType))
	message := map[string]interface{}{
		"event_id":    event.EventID,
		"user_id":     event.UserID,
		"event_type":  event.EventType,
		"payload":     json.RawMessage(event.Payload),
		"occurred_at": event.CreatedAt,
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal lifecycle event: %w", err)
	}
	return p.channel.PublishWithContext(ctx, p.exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

// Close releases the underlying channel and connection.
func (p *RabbitMQPublisher) Close() {
	if p == nil {
		return
	}
	if err := p.channel.Close(); err != nil {
		log.Printf("rabbitmq channel close error: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		log.Printf("rabbitmq connection close error: %v", err)
	}
}
