package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	log.Printf("Received request: Method=%s Path=%s", request.HTTPMethod, request.Path)

	// Get environment variables
	stage := os.Getenv("STAGE")
	logLevel := os.Getenv("LOG_LEVEL")

	log.Printf("Environment: STAGE=%s LOG_LEVEL=%s", stage, logLevel)

	// Parse request body
	var req Request
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Error parsing request: %v", err)
		return Response{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "Invalid request body"}`,
		}, nil
	}

	// Process request
	result := fmt.Sprintf("Hello %s! Your message: %s", req.Name, req.Message)

	// Create response
	responseBody := map[string]interface{}{
		"message": result,
		"stage":   stage,
	}

	body, _ := json.Marshal(responseBody)

	return Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
