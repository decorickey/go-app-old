package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go-app/config"
	"go-app/utils"
	"io/ioutil"
	"net/http"
)

func ApiHandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// リクエスト解析
	requestMethod := request.HTTPMethod
	requestPath := request.Path

	// レスポンス情報
	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}

	if requestMethod == "GET" && requestPath == "/bmonster" {
		programList, err := bmonsterLatestProgramList()
		if err != nil || programList == "" {
			headers["Content-Type"] = "text/plain"
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				Body:       "Internal server error",
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		headers["Content-Type"] = "application/json"
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       programList,
			StatusCode: http.StatusOK,
		}, nil
	}

	headers["Content-Type"] = "text/plain"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       fmt.Sprintf("method: %s, path: %s", requestMethod, requestPath),
		StatusCode: http.StatusBadRequest,
	}, nil
}

func main() {
	lambda.Start(ApiHandleRequest)
}

func bmonsterLatestProgramList() (string, error) {
	// S3からjsonファイル取得
	key := "bmonster.json"
	filename := "/tmp/" + key // Lambdaでは/tmp配下でないと権限がない
	err := utils.GetObject(config.Config.Bucket, key, filename)
	if err != nil {
		return "", err
	}

	// jsonファイル解析
	json, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
