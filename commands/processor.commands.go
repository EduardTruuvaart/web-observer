package commands

import (
	"context"
	"strings"

	"github.com/EduardTruuvaart/web-observer/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgCommandProcessor) processCommands(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	command := extractCommand(message.Text)

	var err error

	switch command {
	case "start":
		_, _ = SendMsg(t.bot, chatID, "Hello there! Just send me full URL that you want to track.")
		err = t.botFlowRepository.Save(ctx, chatID, domain.URLRequsted)
		t.contentRepository.Create(ctx, chatID)

	case "stop":
		_, _ = SendMsg(t.bot, chatID, "Bye bye!")
		err = t.botFlowRepository.Delete(ctx, chatID)
		t.contentRepository.Delete(ctx, chatID)

	default:
		_, err = SendMsg(t.bot, chatID, "Sorry, I don't understand that command. Please use /start to get started.")
	}

	return err
}

// properly extracts command from the input string, removing all unnecessary parts
// please refer to unit tests for details
func extractCommand(rawCommand string) string {
	command := rawCommand

	// remove slash if necessary
	if rawCommand[0] == '/' {
		command = command[1:]
	}

	return strings.TrimSpace(command)
}
