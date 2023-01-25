package main

import (
	"context"
	"log"
	"net/http"

	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	db := *dynamodb.NewFromConfig(cfg)
	httpClient := &http.Client{}
	contentRepository := repository.NewDynamoContentRepository(db)
	contentFetcher := service.NewContentFetcher(contentRepository, httpClient)

	url := "https://eu.store.ui.com/collections/unifi-protect-cameras/products/g4-doorbell-pro"

	_, err = contentFetcher.FetchAndCompare(ctx, url)

	if err != nil {
		log.Fatal(err)
	}
}
