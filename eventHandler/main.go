package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is a function to acting as a webhook for Twitter's Account Activity
// Subscriptions.
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	log.Printf("Request Body: %v", request.Body)

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Event received successfully",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
