package service

import (
	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/models"
)

func CreateUser(clerkUserId string) error {
	var user models.AppUser
	user.ClerkId = clerkUserId

	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}
	go GetOrCreateStripeCustomerIDForUser(user.Id.Bytes)
	return nil
}
