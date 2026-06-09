package model

import (
	"time"
)

const (
	ShortCodeLength = 10
	Charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

type URL struct {
	ShortCode string    `json:"short_code"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
}
