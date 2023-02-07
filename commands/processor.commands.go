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

	// TODO: Use parallel processing
	switch command {
	case "start":
		_, _ = SendMsg(t.bot, chatID, "Hello there! Just send me full URL that you want to track.")
		err = t.botFlowRepository.Save(ctx, chatID, domain.URLRequsted, nil)

	case "stop":
		_, _ = SendMsg(t.bot, chatID, "Bye bye!")
		flowStep, savedUrl, _ := t.botFlowRepository.FindByChatID(ctx, chatID)

		if flowStep != domain.NotStarted {
			err = t.botFlowRepository.Delete(ctx, chatID)

			if savedUrl != nil {
				t.contentRepository.Delete(ctx, chatID, *savedUrl)
			}
		}

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
