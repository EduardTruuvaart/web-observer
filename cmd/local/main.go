package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	//bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	db := *dynamodb.NewFromConfig(cfg)
	s3Client := *s3.NewFromConfig(cfg)

	//proxyUrl, _ := url.Parse("http://localhost:8081")
	httpClient := &http.Client{
		// Transport: &http.Transport{
		// 	Proxy:           http.ProxyURL(proxyUrl),
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// },
	}
	contentRepository := repository.NewDynamoContentRepository(db, s3Client, "ObserverTraces", "web-observer-bucket")
	contentFetcher := service.NewContentFetcher(contentRepository, httpClient)

	result, err := contentFetcher.FetchAndCompare(ctx, 493004756, "https://www.johnlewis.com/miele-pur68w-chimney-cooker-hood-stainless-steel/p3095618",
		"div.Layout_Row__ev3af")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %s\n", result.State)
	fmt.Printf("Diff: %s\n", result.Difference)

	//activeTracks, err := contentRepository.FindAllActive(ctx)

	if err != nil {
		log.Fatal(err)
	}

	// var wg sync.WaitGroup
	// for _, track := range activeTracks {
	// 	wg.Add(1)
	// 	go func(track domain.ObserverTrace) {
	// 		defer wg.Done()

	// 		result, err := contentFetcher.FetchAndCompare(ctx, track.ChatID, *track.URL, *track.CssSelector)

	// 		if err != nil {
	// 			log.Fatal(err)
	// 			return
	// 		}

	// 		msg := fmt.Sprintf("❗️ URL %v has changed", *track.URL)
	// 		commands.SendMsg(bot, track.ChatID, msg)

	// 		if result.State == domain.Updated {
	// 			fmt.Printf("Diff size: %v\n", result.DiffSize)
	// 			fmt.Printf("Difference: %s\n", result.Difference)
	// 		}
	// 	}(track)
	// }
	// wg.Wait()
}
