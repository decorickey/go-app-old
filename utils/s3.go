package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io/ioutil"
	"os"
)

// S3PutObjectAPI 定型文
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// PutFile 定型文
func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func PutObject(bucket, key, filename string) error {
	// クライアント生成
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}

	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		return err
	}
	return nil
}

// S3GetObjectAPI 定型文
type S3GetObjectAPI interface {
	GetObject(ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// GetFile 定型文
func GetFile(c context.Context, api S3GetObjectAPI, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return api.GetObject(c, input)
}

func GetObject(bucket, key, filename string) (err error) {
	// クライアント生成
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	content, err := GetFile(context.TODO(), client, input)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(content.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, body, 0666)
	if err != nil {
		return err
	}
	return nil
}
