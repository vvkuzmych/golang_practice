package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Event represents input to Lambda
type Event struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Response represents Lambda output
type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// ResponseBody for JSON body
type ResponseBody struct {
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	ColdStart bool   `json:"coldStart"`
}

var (
	// Track cold start
	coldStart = true
	// Init timestamp
	initTime time.Time
)

func init() {
	initTime = time.Now()
	log.Printf("Lambda initializing at %v", initTime)

	// Simulate some init work
	// In real app: connect to DB, load config, etc.
	time.Sleep(100 * time.Millisecond)

	log.Printf("Init completed in %v", time.Since(initTime))
}

// handler is the Lambda function handler
func handler(ctx context.Context, event Event) (Response, error) {
	// Get Lambda context
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return Response{}, fmt.Errorf("failed to get lambda context")
	}

	// Log cold/warm start
	startType := "WARM"
	if coldStart {
		startType = "COLD"
		log.Printf("COLD START - RequestID: %s", lc.RequestID)
		log.Printf("Time since init: %v", time.Since(initTime))
		coldStart = false
	} else {
		log.Printf("WARM START - RequestID: %s", lc.RequestID)
	}

	// Log incoming event
	log.Printf("Received event: %+v", event)

	// Business logic
	responseMsg := fmt.Sprintf("Hello %s! Your message: %s", event.Name, event.Message)

	// Create response body
	responseBody := ResponseBody{
		Message:   responseMsg,
		RequestID: lc.RequestID,
		Timestamp: time.Now().Format(time.RFC3339),
		ColdStart: startType == "COLD",
	}

	bodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("ERROR: Failed to marshal response: %v", err)
		return Response{}, err
	}

	// Return response
	response := Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(bodyJSON),
	}

	log.Printf("Returning response: %d", response.StatusCode)
	return response, nil
}

func main() {
	lambda.Start(handler)
}
