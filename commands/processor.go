package commands

import (
	"context"
	"io"
	"strings"

	"github.com/EduardTruuvaart/web-observer/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// CommandProcessor processes a Telegram update message
type CommandProcessor interface {
	io.Closer
	Process(ctx context.Context, update tgbotapi.Update) error
}

type TgCommandProcessor struct {
	bot               *tgbotapi.BotAPI
	contentRepository repository.ContentRepository
	botFlowRepository repository.BotFlowRepository
}

func NewProcessor(contentRepository repository.ContentRepository, botFlowRepository repository.BotFlowRepository,
	bot *tgbotapi.BotAPI) (*TgCommandProcessor, error) {
	return &TgCommandProcessor{
		bot:               bot,
		contentRepository: contentRepository,
		botFlowRepository: botFlowRepository,
	}, nil
}

// Process an event, any update happened in telegram bot will be processed by one entry point
func (t TgCommandProcessor) Process(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {

		if update.Message.IsCommand() {
			// this is a command, i.e. starting with /, for example, "/start" or "/track'
			return t.processCommands(ctx, update.Message)
		} else {
			// this is a plain message, i.e. not starting with /, for example, "https://google.com"
			return t.processMessages(ctx, update.Message)
		}

		// just a plain message was sent
		SendMsg(t.bot, update.Message.Chat.ID,
			"Sorry, I don't understand you. Please use /start to get started.")
		return nil
	}

	if update.CallbackQuery != nil {

		// this is button click callback
		//return t.processButtonCallback(ctx, update.CallbackQuery)
	}

	return nil
}

// SendMsg simply send a message to bot in Markdown format, general function
func SendMsg(bot *tgbotapi.BotAPI, chatID int64, textMarkdown string) (tgbotapi.Message, error) {

	if bot == nil {
		return tgbotapi.Message{}, errors.New("Bot was not set")
	}

	msg := tgbotapi.NewMessage(chatID, textMarkdown)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	// send the message
	resp, err := bot.Send(msg)
	if err != nil {
		return resp, errors.Wrap(err, "Error sending message")
	}

	return resp, nil
}

// From DOCs: https://core.telegram.org/bots/api#markdownv2-style
// esc
// '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'
// must be escaped with the preceding character '\'.`'
func markdownEscape(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(s)
}
