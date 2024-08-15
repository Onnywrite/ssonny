package tokens

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpired        = errors.New("token expired")
	ErrIssuerMismatch = errors.New("issuer mismatch")
)

type Id struct {
	Issuer          string    `json:"iss"`
	Subject         uuid.UUID `json:"sub"`
	Audience        uint64    `json:"aud"`
	AuthorizedParty string    `json:"azp"`
	ExpiresAt       int64     `json:"exp"`

	Nickname   *string    `json:"nickname"`
	Email      string     `json:"email"`
	IsVerified bool       `json:"verified"`
	Gender     *string    `json:"gender"`
	Birthday   *time.Time `json:"birthday"`
	Roles      []string   `json:"roles"`
}

func (a Id) Valid() error {
	return nil
}

type Access struct {
	Issuer          string    `json:"iss"`
	Subject         uuid.UUID `json:"sub"`
	Audience        uint64    `json:"aud"`
	AuthorizedParty string    `json:"azp"`
	ExpiresAt       int64     `json:"exp"`
	Scopes          []string  `json:"scopes"`
}

func (a Access) Valid() error {
	return nil
}

type Refresh struct {
	Issuer          string    `json:"iss"`
	Subject         uuid.UUID `json:"sub"`
	Audience        uint64    `json:"aud"`
	AuthorizedParty string    `json:"azp"`
	ExpiresAt       int64     `json:"exp"`
	Id              uint64    `json:"jid"`
	Rotation        uint64    `json:"rtn"`
}

func (a Refresh) Valid() error {
	return nil
}
