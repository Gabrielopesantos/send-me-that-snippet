package model

import (
	"time"
)

type Paste struct {
	Id         string        `json:"id,omitempty" validate:"isdefault" gorm:"primaryKey"`
	Content    string        `json:"content" validate:"gt=0,lte=300"`
	ContentSha string        `json:"content_sha,omitempty" validate:"isdefault"`
	Language   string        `json:"language"`
	CreatedAt  time.Time     `json:"created_at,omitempty" validate:"isdefault" gorm:"autoCreateTime"`
	ExpiresIn  time.Duration `json:"expires_in"`
	Expired    bool
}
