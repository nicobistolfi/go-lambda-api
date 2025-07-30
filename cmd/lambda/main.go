package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nicobistolfi/go-lambda-api/internal/handlers"
	_ "github.com/nicobistolfi/go-lambda-api/internal/models"

	_ "github.com/nicobistolfi/go-lambda-api/docs"
)

// @title						Go Lambda API
// @version					1.0
// @description				A serverless API built with AWS Lambda and Go
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-API-Key
func handleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	//	@Summary		Health check endpoint
	//	@Description	Check if the API is running
	//	@Tags			health
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{object}	handlers.HealthResponse
	//	@Router			/health [get]
	// Handle health endpoint
	if request.RawPath == "/health" {
		response := handlers.HealthResponse{
			Status: "ok",
		}

		body, err := json.Marshal(response)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       `{"error":"Internal server error"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	//	@Summary		Protected endpoint example
	//	@Description	All endpoints except /health require API key authentication
	//	@Tags			authentication
	//	@Accept			json
	//	@Produce		json
	//	@Security		ApiKeyAuth
	//	@Failure		401	{object}	models.ErrorResponse
	//	@Failure		404	{object}	models.ErrorResponse
	// For other endpoints, check authentication
	apiKey, exists := request.Headers["x-api-key"]
	if !exists {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"error":"Missing API key"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Validate API key
	expectedAPIKey := os.Getenv("API_KEY")
	if expectedAPIKey == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"error":"API key not configured"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	if apiKey != expectedAPIKey {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"error":"Invalid API key"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Handle other routes here in the future

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"error":"Not found"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
