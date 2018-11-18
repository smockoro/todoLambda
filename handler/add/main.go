package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/smockoro/todoLambda/driver/db"
)

type request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortenResource string `json:"shorten_resource"`
}

type Link struct {
	ShortenResource string `json:"shorten_resource"`
	OriginalURL     string `json:"original_url"`
}

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return response(http.StatusOK, b), nil
}

func parseRequest(req events.APIGatewayProxyRequest) (*request, error) {
	if req.HTTPMethod != http.MethodPost {
		return nil, fmt.Errorf("use POST request")
	}
	var r request
	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse request")
	}
	// ParseRequestURI は文字列を受け取り、URL にパースするメソッドです。 // エラーなくパースできることで有効な URL とみなしています。
	// https://golang.org/src/net/url/url.go?s=13616:13665#L471
	_, err = url.ParseRequestURI(r.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid URL")
	}
	return &r, nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	// Lambda プロキシ統合のレスポンスフォーマットに沿った構造体が // aws-lambda-go で定義されています。
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func responseBody(shortenResource string) (string, error) {
	resp, err := json.Marshal(Response{ShortenResource: shortenResource})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}
