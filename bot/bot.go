package bot

import (
	"TeleHook/structs"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gofiber/fiber/v2"
)

var (
	bMapMutex sync.Mutex
	TGBot     = map[string]*tgbotapi.BotAPI{}
)

// Start a new bot
func StartBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bMapMutex.Lock()
	TGBot[token] = bot
	bMapMutex.Unlock()

	return bot, nil
}

func BotHandler(c *fiber.Ctx) error {
	urlArgs := strings.Split(c.OriginalURL(), "/")

	// Check if path is too short to contain bottoken
	if len(urlArgs) < 3 {
		c.SendString("No BotToken Supplied")
		c.SendStatus(400)
		return nil
	}

	botToken := urlArgs[2]

	var botMessage structs.MESSAGE
	jsonErr := json.Unmarshal(c.Body(), &botMessage)
	if jsonErr != nil {
		c.SendString(structs.TH_WEB_JSON_FAIL)
		c.SendStatus(400)
		return nil
	}

	bMapMutex.Lock()
	botInstance := TGBot[botToken]
	bMapMutex.Unlock()

	if botInstance == nil {
		var bStartErr error
		botInstance, bStartErr = StartBot(botToken)
		if bStartErr != nil {
			c.SendString(bStartErr.Error())
			c.SendStatus(400)
			return nil
		}
	}

	sMsgErr := sendMessage(botInstance, botMessage)
	if sMsgErr != nil {
		c.SendString(sMsgErr.Error())
		c.SendStatus(400)
		return nil
	}

	return nil
}

func sendMessage(bot *tgbotapi.BotAPI, message structs.MESSAGE) error {

	// No image is being send
	if message.Image == "" {

		botMessage := tgbotapi.NewMessage(message.Channel, message.Message)

		// Create "embed" buttons
		botMessage.ReplyMarkup = processMarkup(message)

		botMessage.ParseMode = message.Markdown

		_, err := bot.Send(botMessage)
		if err != nil {
			return err
		}
	} else {

		imageData, imgErr := base64.StdEncoding.DecodeString(message.Image)
		if imgErr != nil {
			return errors.New(structs.TH_WEB_IMG_FAIL)
		}

		imageFile := tgbotapi.FileBytes{
			Name:  "img.jpg",
			Bytes: imageData,
		}

		botMessage := tgbotapi.NewPhotoUpload(message.Channel, imageFile)

		botMessage.Caption = message.Message
		botMessage.ParseMode = message.Markdown

		// Create "embed" buttons
		botMessage.ReplyMarkup = processMarkup(message)

		_, err := bot.Send(botMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func processMarkup(message structs.MESSAGE) tgbotapi.InlineKeyboardMarkup {
	var inlineRows [][]tgbotapi.InlineKeyboardButton
	for _, row := range message.EmbedRows {
		var inlineButtons []tgbotapi.InlineKeyboardButton
		for _, embed := range row {
			if embed.Action == "" {
				embed.Action = "0"
			}

			var button tgbotapi.InlineKeyboardButton

			switch embed.Type {
			case 1:
				button = tgbotapi.NewInlineKeyboardButtonURL(embed.Name, embed.Action)
			default:
				button = tgbotapi.NewInlineKeyboardButtonData(embed.Name, embed.Action)
			}

			inlineButtons = append(inlineButtons, button)
		}
		inlineRow := tgbotapi.NewInlineKeyboardRow(inlineButtons...)
		inlineRows = append(inlineRows, inlineRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(inlineRows...)
}
