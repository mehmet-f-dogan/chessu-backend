package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

	"gorm.io/gorm/clause"

	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/models"
	"mehmetfd.dev/chessu-backend/service"
)

func AssignMembershipHandlers(app *fiber.App) {
	app.Post("/membership/:userId/cancel", handleCancelMembership)
	app.Get("/membership/:userId/verify", handleVerifyMembership)
	app.Post("/membership/:userId/create-checkout-link", handleCreateMembershipCheckoutLink)
}

func handleCancelMembership(c *fiber.Ctx) error {
	clerkUserId := utils.CopyString(c.Params("userId"))

	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if user.Membership == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	err := service.CancelMembership(user.Id.Bytes)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func handleVerifyMembership(c *fiber.Ctx) error {
	clerkUserId := utils.CopyString(c.Params("userId"))

	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	if user.Membership == nil {
		return c.JSON(fiber.Map{
			"verified": false,
		})
	}

	return c.JSON(fiber.Map{
		"verified": true,
	})
}

func handleCreateMembershipCheckoutLink(c *fiber.Ctx) error {
	clerkUserId := utils.CopyString(c.Params("userId"))
	url, err := service.GenerateMembershipCheckoutLink(clerkUserId)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}
