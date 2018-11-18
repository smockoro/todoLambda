package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/smockoro/todoLambda/domain"
	"github.com/smockoro/todoLambda/driver/db"
)

type request struct {
	User    string `json:"user"`
	Subject string `json:"subject"`
}

type Response struct {
	Id string `json:"id"`
}

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	t := time.Now()
	converted := sha256.Sum256([]byte(t.String() + p.User + p.Subject))
	id := hex.EncodeToString(converted[:])
	todo := &model.Todo{
		Id:      id,
		User:    p.User,
		Subject: p.Subject,
		Status:  "none",
	}

	_, err = DynamoDB.PutItem(todo)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}

	b, err := responseBody(id)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}
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
	return &r, nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func responseBody(id string) (string, error) {
	resp, err := json.Marshal(Response{Id: id})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}
