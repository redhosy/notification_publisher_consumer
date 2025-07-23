package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"rabbitmq-go-example/pkg/rabbitmq"
	"strings"
	"syscall"
	"time"
)

func main() {
	fmt.Println("RabbitMQ Publisher Starting...")

	// Define RabbitMQ connection parameters
	queueName := "notification_queue"
	exchange := "notification_exchange"
	routingKey := "notification_key"

	// Create RabbitMQ client
	rabbitMQ := rabbitmq.NewRabbitMQ(queueName, exchange, routingKey, "")
	defer rabbitMQ.Close()

	// Print connection info
	rabbitMQ.PrintConnectionInfo()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create input channel
	inputChan := make(chan string)

	// Start goroutine to read from stdin
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter message (or 'exit' to quit): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				continue
			}

			// Trim newline characters
			input = strings.TrimSpace(input)

			if input == "exit" {
				close(inputChan)
				return
			}

			inputChan <- input
		}
	}()

	// Start automatic publisher that sends a message every 10 seconds
	go func() {
		counter := 1
		for {
			msg := fmt.Sprintf("Automatic message #%d - %s", counter, time.Now().Format(time.RFC3339))
			err := rabbitMQ.Publish(msg)
			if err != nil {
				fmt.Printf("Failed to publish automatic message: %v\n", err)
			} else {
				fmt.Printf("Published automatic message: %s\n", msg)
			}
			counter++
			time.Sleep(10 * time.Second)
		}
	}()

	// Main loop
	running := true
	for running {
		select {
		case input, ok := <-inputChan:
			if !ok {
				running = false
				break
			}

			// Publish message
			err := rabbitMQ.Publish(input)
			if err != nil {
				fmt.Printf("Failed to publish message: %v\n", err)
			} else {
				fmt.Printf("Published message: %s\n", input)
			}

		case sig := <-sigChan:
			fmt.Printf("Received signal: %v. Shutting down...\n", sig)
			running = false
		}
	}

	fmt.Println("Publisher stopped.")
}
