package main

import (
	"fmt"
	"os"
	"os/signal"
	"rabbitmq-go-example/pkg/rabbitmq"
	"strings"
	"syscall"
	"time"
)

func main() {
	fmt.Println("RabbitMQ Consumer Starting...")

	// Define RabbitMQ connection parameters - must match publisher
	queueName := "notification_queue"
	exchange := "notification_exchange"
	routingKey := "notification_key"

	// Create RabbitMQ client
	rabbitMQ := rabbitmq.NewRabbitMQ(queueName, exchange, routingKey, "")
	defer rabbitMQ.Close()

	// Print connection info
	rabbitMQ.PrintConnectionInfo()

	// Start consuming messages
	msgs, err := rabbitMQ.Consume()
	if err != nil {
		fmt.Printf("Failed to register a consumer: %v\n", err)
		return
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	go func() {
		for msg := range msgs {
			// Process the message
			receivedTime := time.Now().Format(time.RFC3339)
			fmt.Printf("[%s] Received message: %s\n", receivedTime, string(msg.Body))

			// Simulate processing time
			time.Sleep(500 * time.Millisecond)

			// Acknowledge the message
			msg.Ack(false)
			fmt.Println("Message processed and acknowledged.")

			// Add a separator for readability
			fmt.Println(strings.Repeat("-", 50))
		}
	}()

	fmt.Println("Waiting for messages. Press CTRL+C to exit.")

	// Wait for termination signal
	<-sigChan
	fmt.Println("Shutting down consumer...")
}
