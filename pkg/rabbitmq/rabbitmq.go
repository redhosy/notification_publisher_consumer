package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ wraps the connection and channel to the RabbitMQ server
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// Queue details
	QueueName  string
	Exchange   string
	RoutingKey string
	// Connection info
	URI string
}

// NewRabbitMQ creates a new RabbitMQ client
func NewRabbitMQ(queueName, exchange, routingKey, uri string) *RabbitMQ {
	if uri == "" {
		uri = "amqp://guest:guest@localhost:5672/"
	}

	rabbitMQ := &RabbitMQ{
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routingKey,
		URI:        uri,
	}

	var err error

	// Connect to RabbitMQ server
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.URI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Create a channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare an exchange
	if exchange != "" {
		err = rabbitMQ.channel.ExchangeDeclare(
			exchange, // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare an exchange: %v", err)
		}
	}

	// Declare a queue
	_, err = rabbitMQ.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind the queue to the exchange
	if exchange != "" {
		err = rabbitMQ.channel.QueueBind(
			queueName,  // queue name
			routingKey, // routing key
			exchange,   // exchange
			false,      // no-wait
			nil,        // arguments
		)
		if err != nil {
			log.Fatalf("Failed to bind a queue: %v", err)
		}
	}

	return rabbitMQ
}

// Publish sends a message to the queue
func (r *RabbitMQ) Publish(message string) error {
	return r.channel.Publish(
		r.Exchange,   // exchange
		r.RoutingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		})
}

// Consume registers a consumer and handles messages
func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		r.QueueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}

// Close closes the connection to RabbitMQ
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// PrintConnectionInfo prints the connection information
func (r *RabbitMQ) PrintConnectionInfo() {
	fmt.Printf("Connected to RabbitMQ at %s\n", r.URI)
	fmt.Printf("Using queue: %s, exchange: %s, routing key: %s\n",
		r.QueueName, r.Exchange, r.RoutingKey)
}
