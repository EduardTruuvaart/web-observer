package main

import (
	"context"
	"log"
	"os"

	"github.com/EduardTruuvaart/web-observer/repository"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	chatID := int64(493004756)
	msg := tgbotapi.NewMessage(chatID, "Hello World")
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}
