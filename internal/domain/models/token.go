package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Id        uint64    `db:"token_id"`
	UserId    uuid.UUID `db:"token_user_fk"`
	AppId     uint64    `db:"token_app_fk"`
	Rotation  uint64    `db:"token_rotation"`
	RotatedAt time.Time `db:"token_rotated_at"`

	Platform string `db:"token_platform"`
	Agent    string `db:"token_agent"`
}
