package main

import (
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
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

// DeleteWebhookHandler ...
func DeleteWebhookHandler(request events.APIGatewayProxyRequest) (Response, error) {

	httpClient := createClient()

	webhookID := request.QueryStringParameters["webhook_id"]
	if webhookID == "" {
		// return bad request
		return Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            "No webhook_id param present",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	path := fmt.Sprintf("https://api.twitter.com/1.1/account_activity/all/dev/webhooks/%s.json", webhookID)
	req, _ := http.NewRequest("DELETE", path, nil)

	resp, err := httpClient.Do(req)
	if err != nil || resp.StatusCode != 204 {
		return Response{
			StatusCode: resp.StatusCode,
			Body:       err.Error(),
		}, err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 204 {
		return Response{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}, err
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Webhook deleted successfully",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(DeleteWebhookHandler)
}
