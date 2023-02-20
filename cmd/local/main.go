package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	db := *dynamodb.NewFromConfig(cfg)
	s3Client := *s3.NewFromConfig(cfg)
	httpClient := &http.Client{}
	contentRepository := repository.NewDynamoContentRepository(db, s3Client, "ObserverTraces", "web-observer-bucket")
	contentFetcher := service.NewContentFetcher(contentRepository, httpClient)

	activeTracks, err := contentRepository.FindAllActive(ctx)

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, track := range activeTracks {
		wg.Add(1)
		go func(track domain.ObserverTrace) {
			defer wg.Done()

			result, err := contentFetcher.FetchAndCompare(ctx, track.ChatID, *track.URL, *track.CssSelector)

			if err != nil {
				log.Fatal(err)
				return
			}

			if result.State == domain.Updated {
				fmt.Printf("Diff size: %v\n", result.DiffSize)
				fmt.Printf("Difference: %s\n", result.Difference)
			}
		}(track)
	}
	wg.Wait()
}
