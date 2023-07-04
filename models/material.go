package models

import "mehmetfd.dev/chessu-backend/lib"

type Course struct {
	Id            lib.UUID  `json:"id"`
	Chapters      []Chapter `json:"chapters"`
	StripePriceId string    `json:"stripePriceId"`
}

type Chapter struct {
	Id       lib.UUID  `json:"id"`
	Contents []Content `json:"contents"`
	IsSample bool      `json:"isSample"`
}

type Content struct {
	Id lib.UUID `json:"id"`
}
