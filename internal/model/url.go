package model

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID        uuid.UUID `json:"id"`
	Original  string    `json:"original"`
	ShortCode string    `json:"short_code"`
	CreatedAt time.Time `json:"created_at"`
}

func (u URL) Clone() *URL {
	copied := u
	return &copied
}
