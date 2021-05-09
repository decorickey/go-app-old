package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"go-app/app/models"
	"go-app/bmonster"
	"go-app/config"
	"go-app/utils"
	"io/ioutil"
	"log"
)

type Event struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	if event.Name == "bmonster" {
		if err := bmonsterScraping(); err != nil {
			log.Println(err)
			return "Internal server error", err
		}
		return "Success", nil
	}
	return "No such event.", nil
}

func main() {
	lambda.Start(HandleRequest)
}

func bmonsterScraping() error {
	// DynamoDBからスタジオ一覧取得
	studioList, err := models.ScanAllStudio()
	if err != nil {
		return err
	}

	// スクレイピング
	apiClient := bmonster.New("")
	ch := make(chan []map[string]string)
	for _, studio := range studioList {
		go apiClient.ScrapingProgramList(ch, studio)
	}
	programList := make([]map[string]string, 0)
	for i := 0; i < len(studioList); i++ {
		p := <- ch
		programList = append(programList, p...)
	}

	// S3にアップロード
	js, err := json.Marshal(programList)
	if err != nil {
		return err
	}
	key := "bmonster.json"
	filename := "/tmp/" + key // Lambdaでは/tmp配下でないと権限がない
	if err = ioutil.WriteFile(filename, js, 0666); err != nil {
		return err
	}
	err = utils.PutObject(config.Config.Bucket, key, filename)
	if err != nil {
		return err
	}
	return nil
}
