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

// GetWebhooksHandler ...
func GetWebhooksHandler(request events.APIGatewayProxyRequest) (Response, error) {

	httpClient := createClient()

	path := `https://api.twitter.com/1.1/account_activity/all/webhooks.json`
	req, _ := http.NewRequest("GET", path, nil)

	bearer := fmt.Sprintf("bearer %s", os.Getenv("TWITTER_BEARER_TOKEN"))
	req.Header.Add("Authorization", bearer)

	resp, err := httpClient.Do(req)
	if err != nil {
		return Response{
			StatusCode: resp.StatusCode,
			Body:       err.Error(),
		}, err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(respBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(GetWebhooksHandler)
}
