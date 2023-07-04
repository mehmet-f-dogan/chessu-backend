package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"mehmetfd.dev/chessu-backend/controller"
	"mehmetfd.dev/chessu-backend/controller/webhook"
	"mehmetfd.dev/chessu-backend/service"

	"mehmetfd.dev/chessu-backend/database"
)

var environmentVariables = []string{
	"POSTGRES_HOST",
	"POSTGRES_PORT",
	"POSTGRES_USER",
	"POSTGRES_PASSWORD",
	"POSTGRES_DB",

	"REDIS_HOST",
	"REDIS_PORT",

	"AWS_KEY",
	"AWS_SECRET",
	"AWS_MATERIALS_S3_BUCKET_NAME",

	"CLERK_WEBHOOK_SECRET",

	"STRIPE_SECRET_KEY",
	"STRIPE_WEBHOOK_SECRET",
	"MEMBERSHIP_STRIPE_PRICE_ID",
	"MEMBERSHIP_COUPON_ID",

	"APPLICATION_PORT",
}

func main() {
	err := godotenv.Load(".env")
	println("Loaded .env")
	if err != nil {
		panic(err)
	}
	verifyEnvironmentVariables()
	database.InitDB()
	database.LoadMaterials()
	service.InitStripe()

	app := fiber.New()

	controller.AssignMembershipHandlers(app)
	controller.AssignCompletionHandlers(app)

	controller.AssignCoursePurchaseHandlers(app)

	controller.AssignHomepageHandlers(app)

	controller.AssignMembershipHandlers(app)
	webhook.AssignWebhookHandlers(app)

	port := os.Getenv("APPLICATION_PORT")

	err = app.Listen(fmt.Sprintf(":%s", port))
	panic(err)
}

func verifyEnvironmentVariables() {
	for _, envVar := range environmentVariables {
		value := os.Getenv(envVar)
		if value == "" {
			panic("Missing environment variable: " + envVar)
		}
	}
}
