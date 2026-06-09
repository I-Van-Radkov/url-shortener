package model

import (
	"time"
)

type URL struct {
	ShortCode string    `json:"short_code"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
}

func (u URL) Clone() *URL {
	copied := u
	return &copied
}
