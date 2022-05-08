package model

import (
	"time"
)

type Paste struct {
	Id         string        `json:"id,omitempty"`
	Content    string        `json:"content" validate:"gt=0,lte=300"`
	ContentSha string        `json:"content_sha,omitempty"`
	Language   string        `json:"language"`
	CreatedAt  time.Time     `json:"created_at,omitempty"`
	ExpiresIn  time.Duration `json:"expires_in"`
}
