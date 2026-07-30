package model

import (
	"time"

	"github.com/google/uuid"
)

type Cliente struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"passwordHash"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
