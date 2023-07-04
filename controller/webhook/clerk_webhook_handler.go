package webhook

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	svix "github.com/svix/svix-webhooks/go"
	"mehmetfd.dev/chessu-backend/service"
)

var (
	clerkWebhookSecret string
)

func InitClerkWebhookHandler() {
	clerkWebhookSecret = os.Getenv("CLERK_WEBHOOK_SECRET")

}

type ClerkWebhookPayload struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func HandleClerkWebhook(c *fiber.Ctx) error {
	clerkSvixId := c.Get("svix-id")
	clerkSvixTimestamp := c.Get("svix-timestamp")
	clerkSvixSignature := c.Get("svix-signature")

	body := c.Body()

	headers := http.Header{}
	headers.Set("svix-id", clerkSvixId)
	headers.Set("svix-timestamp", clerkSvixTimestamp)
	headers.Set("svix-signature", clerkSvixSignature)

	wh, err := svix.NewWebhook(clerkWebhookSecret)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = wh.Verify(body, headers)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var payload ClerkWebhookPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if payload.Type != "user.created" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	clerkUserId, ok := payload.Data["id"].(string)
	if !ok {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = service.CreateUser(clerkUserId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
