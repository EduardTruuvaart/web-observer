package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var contentRepository *repository.DynamoContentRepository

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	db := *dynamodb.NewFromConfig(cfg)
	s3Client := *s3.NewFromConfig(cfg)
	contentRepository = repository.NewDynamoContentRepository(db, s3Client, os.Getenv("TABLE_NAME"), os.Getenv("BUCKET_NAME"))
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context) error {
	httpClient := &http.Client{}
	contentFetcher := service.NewContentFetcher(contentRepository, httpClient)

	url := "https://eu.store.ui.com/collections/unifi-protect-cameras/products/g4-doorbell-pro"

	result, err := contentFetcher.FetchAndCompare(ctx, url, "div.comProduct__title--wrapper > div > span")

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Result: %v\n", result.State)
	if result.State == domain.Updated {
		fmt.Printf("Diff size: %v\n", result.DiffSize)
		fmt.Printf("Difference: %s\n", result.Difference)
	}

	return nil
}
