package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

// WebhookRegistrationHandler ...
func WebhookRegistrationHandler(ctx context.Context) (Response, error) {

	httpClient := createClient()

	//Set parameters
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("TWITTER_WEBHOOK_ENV") + "/webhooks.json"
	values := url.Values{}
	values.Set("url", os.Getenv("BASE_URL")+"/webhook/twitter")

	//Make Oauth Post with parameters
	resp, err := httpClient.PostForm(path, values)
	if err != nil {
		return Response{
			StatusCode: resp.StatusCode,
			Body:       fmt.Sprintf("Error registering webhook: %v", err),
		}, err
	}
	defer resp.Body.Close()

	//Parse response and check response
	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		panic(err)
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Webhook id of " + data["id"].(string) + " has been registered",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(WebhookRegistrationHandler)
}
