package models

import (
	"time"

	"mehmetfd.dev/chessu-backend/lib"
)

type AppUser struct {
	Id                 lib.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ClerkId            string        `gorm:"unique"`
	StripeId           string        `gorm:"unique"`
	CompletedContentId lib.UUIDArray `gorm:"type:uuid[];default:'{}'"`
	PurchasedCourseId  lib.UUIDArray `gorm:"type:uuid[];default:'{}'"`
	Membership         *Membership   `gorm:"foreignKey:UserID"`
}

type Membership struct {
	Id                   lib.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID               lib.UUID `gorm:"type:uuid"`
	StripeSubscriptionID string   `gorm:"type:text"`
	ValidUntil           time.Time
}
