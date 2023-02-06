package commands

import (
	"context"
	"log"

	"github.com/EduardTruuvaart/web-observer/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgCommandProcessor) processMessages(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	flowStep, err := t.botFlowRepository.FindByChatID(ctx, chatID)

	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	// TODO: Use parallel processing
	switch flowStep {
	case domain.URLRequsted:
		t.contentRepository.UpdateWithUrl(ctx, chatID, message.Text)
		_ = t.botFlowRepository.Save(ctx, chatID, domain.SelectorRequested)
		_, err = SendMsg(t.bot, chatID, "Please provide css selector for the element you want to track. For example: div.comProduct--wrapper > div > span")
	case domain.SelectorRequested:
		t.contentRepository.UpdateWithSelectorAndActivate(ctx, chatID, message.Text)
		_ = t.botFlowRepository.Save(ctx, chatID, domain.SelectorRequested)
		_, err = SendMsg(t.bot, chatID, "Tracking has been activated. You will receive notifications once selection changes.")
	case domain.NotStarted:
		_, err = SendMsg(t.bot, chatID, "Please use /start command to start a tracking flow")
	}

	return err
}
