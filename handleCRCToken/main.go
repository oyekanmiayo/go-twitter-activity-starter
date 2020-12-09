package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// CRCTokenHandler will handle the intermittent CRC challenge sent by twitter
func CRCTokenHandler(request events.APIGatewayProxyRequest) (Response, error) {

	crcToken := request.QueryStringParameters["crc_token"]
	if crcToken == "" {
		// return bad request
		return Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            "No crc_token param present",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	//Encrypt and encode in base 64 then return
	h := hmac.New(sha256.New, []byte(os.Getenv("TWITTER_CONSUMER_SECRET")))
	h.Write([]byte(crcToken))

	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//Generate response string map
	response := make(map[string]string)
	response["response_token"] = "sha256=" + encoded

	//Turn response map to json and send it to the writer
	responseJSON, _ := json.Marshal(response)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(responseJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(CRCTokenHandler)
}
