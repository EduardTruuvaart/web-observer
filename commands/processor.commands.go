package commands

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgCommandProcessor) processCommands(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	command := extractCommand(message.Text)

	var err error

	switch command {
	case "start":
		_, err = SendMsg(t.bot, chatID, "Hello there! Just send me your location and I will find nearby stations!")

	case "stop":
		_, err = SendMsg(t.bot, chatID, "Buy buy!")

	default:
		_, err = SendMsg(t.bot, chatID, "Sorry, I don't understand that command. Try /help")
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
