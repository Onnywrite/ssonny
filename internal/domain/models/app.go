package models

import (
	"time"

	"github.com/google/uuid"
)

type App struct {
	Id          uint64    `db:"app_id"`
	OwnerId     uuid.UUID `db:"app_owner_fk"`
	Name        string    `db:"app_name"`
	Description string    `db:"app_description"`
	SecretHash  string    `db:"app_secret_hash"`

	CreatedAt time.Time  `db:"app_created_at"`
	UpdatedAt time.Time  `db:"app_updated_at"`
	DeletedAt *time.Time `db:"app_deleted_at"`
}
