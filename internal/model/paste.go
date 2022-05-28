package model

import (
	"time"
)

type Paste struct {
	Id         string        `json:"id,omitempty" validate:"isdefault" gorm:"primaryKey"`
	Content    string        `json:"content" validate:"gt=0,lte=300" gorm:"not null"`
	ContentSha string        `json:"content_sha,omitempty" validate:"isdefault" gorm:"not null"`
	Language   string        `json:"language" gorm:"not null"`
	CreatedAt  time.Time     `json:"created_at,omitempty" validate:"isdefault" gorm:"autoCreateTime"`
	ExpiresIn  time.Duration `json:"expires_in" gorm:"not null"`
	Expired    bool          `json:"expired" validate:"isdefault" gorm:"not null"`
}
