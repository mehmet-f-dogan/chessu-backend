package webhook

import "github.com/gofiber/fiber/v2"

func AssignWebhookHandlers(app *fiber.App) {
	InitClerkWebhookHandler()
	InitStripeWebhookHandler()
	app.Post("/webhook/stripe", HandleStripeWebhook)
	app.Post("/webhook/clerk", HandleClerkWebhook)
}
