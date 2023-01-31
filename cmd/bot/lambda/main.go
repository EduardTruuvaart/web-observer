package main

import (
	"context"
	"log"
	"os"

	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
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

func handleRequest(c context.Context, request events.APIGatewayProxyRequest) error {
	// httpClient := &http.Client{}
	// contentFetcher := service.NewContentFetcher(contentRepository, httpClient)

	bot, err := tgbotapi.NewBotAPI("TELEGRAM_BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	chatID := int64(000000000)
	msg := tgbotapi.NewMessage(chatID, "Hello World")
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
