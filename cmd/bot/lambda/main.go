package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/EduardTruuvaart/web-observer/commands"
	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var contentRepository *repository.DynamoContentRepository
var botFlowRepository *repository.DynamoBotFlowRepository

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	db := *dynamodb.NewFromConfig(cfg)
	s3Client := *s3.NewFromConfig(cfg)
	contentRepository = repository.NewDynamoContentRepository(db, s3Client, os.Getenv("CONTENT_TABLE_NAME"), os.Getenv("BUCKET_NAME"))
	botFlowRepository = repository.NewDynamoBotFlowRepository(db, os.Getenv("FLOW_TABLE_NAME"))
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) error {
	// httpClient := &http.Client{}
	// contentFetcher := service.NewContentFetcher(contentRepository, httpClient)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return err
	}

	commandProcessor, err := commands.NewProcessor(contentRepository, bot)

	if err != nil {
		return err
	}

	var update tgbotapi.Update
	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		return err
	}

	commandProcessor.Process(ctx, update)

	return nil
}
