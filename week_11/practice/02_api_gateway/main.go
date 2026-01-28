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

var (
	// Environment variables
	stage = os.Getenv("STAGE")
)

func init() {
	log.Printf("Initializing Lambda for stage: %s", stage)
}

// handler handles API Gateway proxy requests
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received %s request to %s", request.HTTPMethod, request.Path)
	log.Printf("Query params: %v", request.QueryStringParameters)
	log.Printf("Headers: %v", request.Headers)

	// Route based on method and path
	switch {
	case request.HTTPMethod == "GET" && request.Path == "/health":
		return handleHealth(ctx, request)
	case request.HTTPMethod == "GET" && request.Path == "/users":
		return handleGetUsers(ctx, request)
	case request.HTTPMethod == "POST" && request.Path == "/users":
		return handleCreateUser(ctx, request)
	case request.HTTPMethod == "GET" && request.PathParameters["id"] != "":
		return handleGetUser(ctx, request)
	default:
		return notFound()
	}
}

// handleHealth returns health check response
func handleHealth(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := map[string]interface{}{
		"status": "healthy",
		"stage":  stage,
	}

	return jsonResponse(200, body)
}

// handleGetUsers returns list of users
func handleGetUsers(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// In real app: fetch from DB
	users := []map[string]interface{}{
		{"id": "1", "name": "Alice", "email": "alice@example.com"},
		{"id": "2", "name": "Bob", "email": "bob@example.com"},
	}

	// Filter by query param
	if name := request.QueryStringParameters["name"]; name != "" {
		filtered := []map[string]interface{}{}
		for _, user := range users {
			if user["name"] == name {
				filtered = append(filtered, user)
			}
		}
		users = filtered
	}

	return jsonResponse(200, users)
}

// handleGetUser returns single user
func handleGetUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := request.PathParameters["id"]

	log.Printf("Fetching user: %s", userID)

	// In real app: fetch from DB
	user := map[string]interface{}{
		"id":    userID,
		"name":  "Alice",
		"email": "alice@example.com",
	}

	return jsonResponse(200, user)
}

// handleCreateUser creates new user
func handleCreateUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request body
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		log.Printf("ERROR: Failed to parse body: %v", err)
		return jsonResponse(400, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate
	if input.Name == "" || input.Email == "" {
		return jsonResponse(400, map[string]string{
			"error": "name and email are required",
		})
	}

	// In real app: save to DB
	user := map[string]interface{}{
		"id":    "123",
		"name":  input.Name,
		"email": input.Email,
	}

	log.Printf("Created user: %+v", user)

	return jsonResponse(201, user)
}

// notFound returns 404 response
func notFound() (events.APIGatewayProxyResponse, error) {
	return jsonResponse(404, map[string]string{
		"error": "Not Found",
	})
}

// jsonResponse creates API Gateway response with JSON body
func jsonResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Printf("ERROR: Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Internal Server Error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(bodyJSON),
	}, nil
}

func main() {
	lambda.Start(handler)
}
