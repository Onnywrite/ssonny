package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `db:"user_id"`
	Nickname     *string   `db:"user_nickname"`
	Email        string    `db:"user_email"`
	IsVerified   bool      `db:"user_verified"`
	Gender       *string   `db:"user_gender"`
	PasswordHash *string   `db:"user_password_hash"`

	Birthday  *time.Time `db:"user_birthday"`
	CreatedAt time.Time  `db:"user_created_at"`
	UpdatedAt time.Time  `db:"user_updated_at"`
	DeletedAt *time.Time `db:"user_deleted_at"`
}
