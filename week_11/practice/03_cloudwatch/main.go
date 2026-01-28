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
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Structured log entry
type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"requestId,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Duration  int64                  `json:"duration,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

var (
	functionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	coldStart    = true
)

func init() {
	// Configure log format
	log.SetFlags(0) // Remove default timestamp (we add our own)
	logInfo("Lambda initialized", nil)
}

// handler demonstrates structured logging to CloudWatch
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	startTime := time.Now()

	// Get Lambda context
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := lc.RequestID

	// Log cold/warm start
	if coldStart {
		logInfo("Cold start detected", map[string]interface{}{
			"requestId": requestID,
		})
		coldStart = false
	} else {
		logInfo("Warm start", map[string]interface{}{
			"requestId": requestID,
		})
	}

	// Log incoming request
	logInfo("Incoming request", map[string]interface{}{
		"method":      request.HTTPMethod,
		"path":        request.Path,
		"queryParams": request.QueryStringParameters,
		"requestId":   requestID,
	})

	// Simulate business logic with logging
	result, err := processRequest(ctx, request)
	if err != nil {
		logError("Request processing failed", err, map[string]interface{}{
			"requestId": requestID,
			"path":      request.Path,
		})

		return errorResponse(500, "Internal Server Error")
	}

	// Calculate duration
	duration := time.Since(startTime).Milliseconds()

	// Log successful response
	logInfo("Request completed successfully", map[string]interface{}{
		"requestId":  requestID,
		"statusCode": 200,
		"duration":   duration,
	})

	return jsonResponse(200, result)
}

// processRequest simulates business logic with various log levels
func processRequest(ctx context.Context, request events.APIGatewayProxyRequest) (map[string]interface{}, error) {
	// Debug log
	logDebug("Processing request", map[string]interface{}{
		"path":   request.Path,
		"method": request.HTTPMethod,
	})

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	// Warning example
	if request.QueryStringParameters["debug"] == "true" {
		logWarn("Debug mode enabled", map[string]interface{}{
			"path": request.Path,
		})
	}

	// Result
	result := map[string]interface{}{
		"message": "Request processed successfully",
		"data": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"path":      request.Path,
		},
	}

	return result, nil
}

// Structured logging functions

func logDebug(message string, extra map[string]interface{}) {
	logEntry(LogEntry{
		Level:     "DEBUG",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Function:  functionName,
		Extra:     extra,
	})
}

func logInfo(message string, extra map[string]interface{}) {
	logEntry(LogEntry{
		Level:     "INFO",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Function:  functionName,
		Extra:     extra,
	})
}

func logWarn(message string, extra map[string]interface{}) {
	logEntry(LogEntry{
		Level:     "WARN",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Function:  functionName,
		Extra:     extra,
	})
}

func logError(message string, err error, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}
	extra["error"] = err.Error()

	logEntry(LogEntry{
		Level:     "ERROR",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Function:  functionName,
		Extra:     extra,
	})
}

func logEntry(entry LogEntry) {
	// Marshal to JSON for structured logging
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple log
		log.Printf("ERROR: Failed to marshal log entry: %v", err)
		return
	}

	// Print JSON log (CloudWatch will parse it)
	fmt.Println(string(jsonBytes))
}

// Response helpers

func jsonResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return errorResponse(500, "Failed to marshal response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(bodyJSON),
	}, nil
}

func errorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       fmt.Sprintf(`{"error":"%s"}`, message),
	}, nil
}

func main() {
	lambda.Start(handler)
}
