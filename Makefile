.PHONY: build run-publisher run-consumer clean

# Build both publisher and consumer
build:
	go build -o bin/publisher cmd/publisher/main.go
	go build -o bin/consumer cmd/consumer/main.go

# Run the publisher
run-publisher:
	go run cmd/publisher/main.go

# Run the consumer
run-consumer:
	go run cmd/consumer/main.go

# Clean built binaries
clean:
	rm -rf bin/
