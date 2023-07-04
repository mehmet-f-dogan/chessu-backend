package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/google/uuid"

	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/models"
	"mehmetfd.dev/chessu-backend/service"
)

func AssignCoursePurchaseHandlers(app *fiber.App) {
	app.Get("/purchase/course/:courseId/user/:userId/verify", handleCoursePurchaseVerification)
	app.Post("/purchase/course/:courseId/user/:userId/create-checkout-link", handleCreateCourseCheckoutLink)
	app.Get("/price/course/:courseId/user/:userId", handleCoursePrice)
}

func handleCoursePurchaseVerification(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	courseIdString := c.Params("courseId")
	courseId, err := uuid.Parse(courseIdString)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, purchasedCourseId := range user.PurchasedCourseId.Elements {
		if purchasedCourseId.Bytes == courseId {
			return c.JSON(fiber.Map{
				"verified": true,
			})
		}
	}

	return c.JSON(fiber.Map{
		"verified": false,
	})

}

func handleCreateCourseCheckoutLink(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	courseIdString := utils.CopyString(c.Params("courseId"))
	courseId, err := uuid.Parse(courseIdString)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	url, err := service.GenerateCourseCheckoutLink(courseId, clerkUserId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"url": url,
	})
}

func handleCoursePrice(c *fiber.Ctx) error {
	clerkUserId := c.Params("userId")

	courseIdString := utils.CopyString(c.Params("courseId"))
	courseId, err := uuid.Parse(courseIdString)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	price, err := service.GetUserCoursePrice(courseId, clerkUserId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"price": price,
	})
}
