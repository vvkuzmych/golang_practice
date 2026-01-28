package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Order represents an order message
type Order struct {
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	CustomerID string  `json:"customer_id"`
	Timestamp  string  `json:"timestamp"`
}

// In-memory cache for idempotency (in production, use DynamoDB/Redis)
var processedMessages = make(map[string]bool)

func handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	log.Printf("Processing batch of %d messages", len(sqsEvent.Records))

	// Track failures for partial batch response
	var batchItemFailures []events.SQSBatchItemFailure

	for _, record := range sqsEvent.Records {
		messageID := record.MessageId
		log.Printf("Processing message: %s", messageID)

		// Idempotency check
		if processedMessages[messageID] {
			log.Printf("Skipping duplicate message: %s", messageID)
			continue
		}

		// Process message
		if err := processMessage(record); err != nil {
			log.Printf("Error processing message %s: %v", messageID, err)

			// Add to failures (will retry, then DLQ after maxReceiveCount)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: messageID,
			})
			continue
		}

		// Mark as processed
		processedMessages[messageID] = true
		log.Printf("Successfully processed message: %s", messageID)
	}

	// Return partial batch response
	return events.SQSEventResponse{
		BatchItemFailures: batchItemFailures,
	}, nil
}

func processMessage(record events.SQSMessage) error {
	// Parse message body
	var order Order
	if err := json.Unmarshal([]byte(record.Body), &order); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate order
	if order.OrderID == "" {
		log.Printf("Missing order_id")
		return fmt.Errorf("missing order_id")
	}

	if order.Amount <= 0 {
		log.Printf("Invalid amount: %f", order.Amount)
		return fmt.Errorf("invalid amount: %f", order.Amount)
	}

	// Simulate processing
	log.Printf("Processing order: ID=%s Customer=%s Amount=$%.2f",
		order.OrderID, order.CustomerID, order.Amount)

	// Simulate work
	time.Sleep(100 * time.Millisecond)

	// Business logic here:
	// - Save to database
	// - Call external API
	// - Send notification
	// - etc.

	log.Printf("Order processed successfully: %s", order.OrderID)
	return nil
}

func init() {
	// Set log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Log environment
	env := os.Getenv("ENVIRONMENT")
	queueURL := os.Getenv("QUEUE_URL")
	dlqURL := os.Getenv("DLQ_URL")

	log.Printf("Environment: %s", env)
	log.Printf("Queue URL: %s", queueURL)
	log.Printf("DLQ URL: %s", dlqURL)
}

func main() {
	lambda.Start(handler)
}
