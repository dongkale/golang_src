package entity

import (
	"time"
)

type Human struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" validate:"required" db:"name"`
	Race      string    `json:"race" validate:"required" db:"race"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
