package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	ShortCodeLength = 10
	Charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

type URL struct {
	ID        uuid.UUID `json:"id"`
	ShortCode string    `json:"short_code"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
}
