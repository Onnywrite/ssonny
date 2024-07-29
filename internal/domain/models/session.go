package models

import "github.com/google/uuid"

type Session struct {
	Id           uuid.UUID `db:"session_id"`
	UserId       uuid.UUID `db:"session_user_fk"`
	AppId        uint64    `db:"session_app_fk"` // 0 treats as nil
	LastRotation uint64    `db:"session_last_rotation"`

	Platform string `db:"session_platform"`
	Agent    string `db:"session_agent"`
}
