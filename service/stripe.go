package service

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/subscription"

	"gorm.io/gorm/clause"

	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/models"
)

const membershipDiscountAmount = 0.5

var (
	stripeSecretKey   string
	frontendURL       string
	membershipPriceID string
	membershipCoupon  string
)

func InitStripe() {
	stripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
	frontendURL = os.Getenv("FRONTEND_SERVER_URL")
	membershipPriceID = os.Getenv("MEMBERSHIP_STRIPE_PRICE_ID")
	membershipCoupon = os.Getenv("MEMBERSHIP_COUPON_ID")
	stripe.Key = stripeSecretKey
}

func GetOrCreateStripeCustomerIDForUser(userId uuid.UUID) (string, error) {
	var user models.AppUser
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return "", err
	}

	if user.StripeId != "" {
		return user.StripeId, nil
	}

	params := &stripe.CustomerParams{}
	params.Metadata = map[string]string{
		"userId": userId.String(),
	}

	customer, err := customer.New(params)
	if err != nil {
		return "", err
	}

	err = database.DB.Model(&user).Update("stripe_id", customer.ID).Error
	if err != nil {
		return "", err
	}

	return customer.ID, nil
}

func GenerateCourseCheckoutLink(courseID uuid.UUID, userID string) (string, error) {
	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: userID}).First(&user).Error; err != nil {
		return "", err
	}

	// Check if the user is the owner of the course
	for _, purchasedCourseID := range user.PurchasedCourseId.Elements {
		if purchasedCourseID.Bytes == courseID {
			return "", errors.New("user already owns this course")
		}
	}

	coursePtr := database.GetCourse(courseID)
	if coursePtr == nil {
		return "", errors.New("course not found")
	}

	lineItems := []*stripe.CheckoutSessionLineItemParams{
		{
			Price:    stripe.String(coursePtr.StripePriceId),
			Quantity: stripe.Int64(1),
		},
	}

	userIdStr, _ := uuid.FromBytes(user.Id.Bytes[:])
	metadata := map[string]string{
		"userId":   userIdStr.String(),
		"courseId": courseID.String(),
		"type":     "course",
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(frontendURL + "/course/" + courseID.String() + "/payment-successful"),
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		Customer:   stripe.String(user.StripeId),
	}

	if user.Membership != nil {
		params.Discounts = []*stripe.CheckoutSessionDiscountParams{
			{
				Coupon: stripe.String(membershipCoupon),
			},
		}
	}

	params.Metadata = metadata

	session, err := session.New(params)

	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func GenerateMembershipCheckoutLink(userID string) (string, error) {
	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where(&models.AppUser{ClerkId: userID}).First(&user).Error; err != nil {
		return "", err
	}

	if user.Membership != nil {
		return "", errors.New("user already has a membership")
	}

	lineItems := []*stripe.CheckoutSessionLineItemParams{
		{
			Price:    stripe.String(membershipPriceID),
			Quantity: stripe.Int64(1),
		},
	}

	userIdStr, _ := uuid.FromBytes(user.Id.Bytes[:])
	metadata := map[string]string{
		"userId": userIdStr.String(),
		"type":   "membership",
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(frontendURL + "/membership/payment-successful"),
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer:   stripe.String(user.StripeId),
	}
	params.Metadata = metadata

	session, err := session.New(params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func CancelMembership(userID uuid.UUID) error {
	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	if user.Membership == nil {
		return errors.New("user does not have a membership")
	}

	existingSubscription, err := subscription.Get(user.Membership.StripeSubscriptionID, nil)
	if err != nil {
		return err
	}

	_, err = subscription.Cancel(existingSubscription.ID, nil)
	if err != nil {
		return err
	}

	membership := user.Membership
	user.Membership = nil

	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}

	if err := database.DB.Delete(membership).Error; err != nil {
		return err
	}

	return nil
}

func GetUserCoursePrice(courseID uuid.UUID, userID string) (float64, error) {
	price, err := getCoursePrice(courseID)
	if err != nil {
		return 0, err
	}

	var user models.AppUser
	if err := database.DB.Preload(clause.Associations).Where(&models.AppUser{ClerkId: userID}).First(&user).Error; err != nil {
		return 0, err
	}

	if user.Membership != nil {
		price *= 1 - membershipDiscountAmount
	}

	return price, nil
}

func getCoursePrice(courseID uuid.UUID) (float64, error) {
	coursePtr := database.GetCourse(courseID)

	if coursePtr == nil {
		return 0, errors.New("course not found")
	}

	price, err := price.Get(coursePtr.StripePriceId, nil)
	if err != nil {
		return 0, err
	}

	return price.UnitAmountDecimal, nil
}
