package model

import (
	"time"
)

const (
	ShortCodeLength = 10
)

type URL struct {
	ShortCode string    `json:"short_code"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
}
