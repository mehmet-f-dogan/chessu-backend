package webhook

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"

	"gorm.io/gorm/clause"

	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/lib"
	"mehmetfd.dev/chessu-backend/models"
)

var (
	stripeWebhookSecret string
)

func InitStripeWebhookHandler() {
	stripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
}

func HandleStripeWebhook(c *fiber.Ctx) error {
	signature := c.Get("Stripe-Signature")
	payload := c.Request().Body()
	return handleStripeWebhookPayload(c, signature, payload)
}

func handleStripeWebhookPayload(c *fiber.Ctx, signature string, payload []byte) error {
	event, err := webhook.ConstructEvent(payload, signature, stripeWebhookSecret)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch event.Type {
	case "checkout.session.completed":
		sessionObj := &stripe.CheckoutSession{}
		err := json.Unmarshal(event.Data.Raw, sessionObj)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		checkoutType := sessionObj.Metadata["type"]
		userId := sessionObj.Metadata["userId"]
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		switch checkoutType {
		case "course":
			courseId, err := uuid.Parse(sessionObj.Metadata["courseId"])
			if err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}
			return handleCoursePurchase(c, courseId, userUUID)
		case "membership":
			validUntil := time.Unix(sessionObj.ExpiresAt, 0).UTC()
			subscriptionId := sessionObj.Subscription.ID
			return handleMembershipPurchase(c, userUUID, subscriptionId, validUntil)
		default:
			return errors.New("invalid checkout type")
		}

	case "invoice.paid":
		invoiceObj := &stripe.Invoice{}
		err := json.Unmarshal(event.Data.Raw, invoiceObj)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		subscription := invoiceObj.Lines.Data[0]
		validUntil := time.Unix(subscription.Period.End, 0).UTC()
		customerId := invoiceObj.Customer.ID
		return handleMembershipRegularPayment(c, customerId, validUntil)

	case "invoice.payment_failed":
		invoiceObj := &stripe.Invoice{}
		err := json.Unmarshal(event.Data.Raw, invoiceObj)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		failedCustomerId := invoiceObj.Customer.ID
		return handleMembershipPaymentFail(c, failedCustomerId)

	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func handleCoursePurchase(c *fiber.Ctx, courseId uuid.UUID, userId uuid.UUID) error {
	coursePtr := database.GetCourse(courseId)
	if coursePtr == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).First(&user, userId).Error; err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	purchasedCourseIds := user.PurchasedCourseId
	purchasedCourseIds.Append(courseId.String())
	user.PurchasedCourseId = purchasedCourseIds
	if err := database.DB.Save(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func handleMembershipPurchase(c *fiber.Ctx, userId uuid.UUID, subscriptionId string, validUntil time.Time) error {
	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).First(&user, userId).Error; err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	membership := user.Membership
	if membership == nil {
		var newMembership models.Membership
		newMembership.ValidUntil = validUntil
		newMembership.StripeSubscriptionID = subscriptionId
		newMembership.UserID = lib.UUID{
			UUID: user.Id.UUID,
		}
		membership = &newMembership // Assign the newMembership to membership pointer
	}

	membershipValidUntil := membership.ValidUntil
	if membershipValidUntil.Before(validUntil) {
		membership.ValidUntil = validUntil
		if err := database.DB.Save(membership).Error; err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	user.Membership = membership
	if err := database.DB.Save(&user).Error; err != nil { // Save the pointer to user
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func handleMembershipRegularPayment(c *fiber.Ctx, stripeCustomerId string, validUntil time.Time) error {
	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where("stripe_id = ?", stripeCustomerId).First(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	membership := user.Membership
	if membership == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if membership.ValidUntil.Before(validUntil) {
		membership.ValidUntil = validUntil
	}

	if err := database.DB.Save(membership).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user.Membership = membership
	if err := database.DB.Save(user).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func handleMembershipPaymentFail(c *fiber.Ctx, stripeCustomerId string) error {
	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where("stripe_id = ?", stripeCustomerId).First(&user).Error; err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	existingMembership := user.Membership
	user.Membership = nil

	if err := database.DB.Save(user).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := database.DB.Delete(existingMembership).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
