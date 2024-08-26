package models

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	Id         uint64     `db:"domain_id"`
	OwnerId    uuid.UUID  `db:"domain_owner_fk"`
	Name       string     `db:"domain_name"`
	IsVerified bool       `db:"domain_verified"`
	VerifiedAt *time.Time `db:"domain_verified_at"`

	CreatedAt time.Time  `db:"domain_created_at"`
	UpdatedAt time.Time  `db:"domain_updated_at"`
	DeletedAt *time.Time `db:"domain_deleted_at"`
}
