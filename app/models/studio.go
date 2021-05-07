package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Studio struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var tableName = "bmonster_studio"

// DynamoDBScanAPI 定型文
type DynamoDBScanAPI interface {
	Scan(ctx context.Context,
		params *dynamodb.ScanInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// GetItems 定型文
func GetItems(c context.Context, api DynamoDBScanAPI, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return api.Scan(c, input)
}

func ScanAllStudio() (studioList []Studio, err error) {
	// クライアント生成
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	// スキャン設定
	proj := expression.NamesList(expression.Name("name"), expression.Name("code"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &tableName,
	}

	// 実行
	res, err := GetItems(context.TODO(), client, input)
	if err != nil {
		return nil, err
	}

	// 解析
	err = attributevalue.UnmarshalListOfMaps(res.Items, &studioList)
	if err != nil {
		return nil, err
	}
	return studioList, nil
}
