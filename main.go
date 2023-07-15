package main

import (
	"TeleHook/bot"
	"os"

	fiber "github.com/gofiber/fiber/v2"
)

const (
	version    = "v1.0.0"
	secretPath = "Q75k9anIncOQO9peWkF0HMTkIyjQVsSd"
)

func main() {

	webhook := fiber.New()

	webhook.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("TeleHook " + version)
	})

	webhook.Get("/"+secretPath+"/*", bot.BotHandler)
	webhook.Post("/"+secretPath+"/*", bot.BotHandler)

	webhook.Listen(os.Args[1])
}
