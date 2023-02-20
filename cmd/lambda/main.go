package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/EduardTruuvaart/web-observer/commands"
	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/EduardTruuvaart/web-observer/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var contentRepository *repository.DynamoContentRepository
var httpClient *http.Client
var contentFetcher *service.ContentFetcher

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	db := *dynamodb.NewFromConfig(cfg)
	s3Client := *s3.NewFromConfig(cfg)
	contentRepository = repository.NewDynamoContentRepository(db, s3Client, os.Getenv("CONTENT_TABLE_NAME"), os.Getenv("BUCKET_NAME"))
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context) error {
	httpClient = &http.Client{}
	contentFetcher = service.NewContentFetcher(contentRepository, httpClient)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	activeTracks, err := contentRepository.FindAllActive(ctx)

	if err != nil {
		log.Fatal(err)
		return err
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

			fmt.Printf("ChatID: %d, URL: %s Diff Result: %v\n", track.ChatID, *track.URL, result.State)
			if result.State == domain.Updated {
				msg := fmt.Sprintf("❗️ URL %v has changed", *track.URL)
				commands.SendMsg(bot, track.ChatID, msg)

				fmt.Printf("Diff size: %v\n", result.DiffSize)
				fmt.Printf("Difference: %s\n", result.Difference)
			}
		}(track)
	}
	wg.Wait()

	fmt.Printf("Completed checking %d records \n", len(activeTracks))
	return nil
}
