package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/oauth1"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func createClient() *http.Client {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

// WebhookSubscriptionHandler adds a subscription for the registered webhook url
func WebhookSubscriptionHandler(ctx context.Context) (Response, error) {
	client := createClient()
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("TWITTER_WEBHOOK_ENV") + "/subscriptions.json"

	resp, err := client.PostForm(path, nil)
	if err != nil {
		// something
		return Response{
			StatusCode:      resp.StatusCode,
			IsBase64Encoded: false,
			Body:            fmt.Sprintf("Error subscribing webhook: %v", err),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 204 {
		return Response{
			StatusCode:      resp.StatusCode,
			IsBase64Encoded: false,
			Body:            string(body),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Subscription for " + path + " was successful",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, err
}

func main() {
	lambda.Start(WebhookSubscriptionHandler)
}
