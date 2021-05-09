# go-app

## main

### main.go

WebAPIサーバー用main関数

### lambdahandler.go

スクレイピング処理用AWSLambda関数

### apihandler.go

APIGW+LambdaのAPI用Lambdaかんすう

## AWS Lambda

### zipのビルド

```bash
GOOS=linux go build lambdahandler.go && zip lambdahandler.zip lambdahandler config.ini
GOOS=linux go build apihandler.go && zip apihandler.zip apihandler config.ini
```