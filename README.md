# go-app

## AWS Lambda

### zipのビルド

```bash
GOOS=linux go build lambdahandler.go && zip lambdahandler.zip lambdahandler config.ini
GOOS=linux go build apihandler.go && zip apihandler.zip apihandler config.ini
```