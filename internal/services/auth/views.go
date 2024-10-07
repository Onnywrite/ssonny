package auth

import (
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"

	"github.com/google/uuid"
)

type UserInfo struct {
	Platform string `validate:"max=255,ascii"`
	Agent    string `validate:"max=255,ascii"`
}

type RegisterWithPasswordData struct {
	Nickname *string    `validate:"omitempty,min=3,max=32,ascii"`
	Email    string     `validate:"email,max=345"`
	Gender   *string    `validate:"omitempty,max=16"`
	Birthday *time.Time `validate:"-"`
	Password string     `validate:"min=8,max=72"`
	UserInfo UserInfo
}

type AuthenticatedUser struct {
	Access  string
	Refresh string
	Profile Profile
}

type LoginWithPasswordData struct {
	Email    *string `validate:"omitempty,email,max=345"`
	Nickname *string `validate:"omitempty,min=3,max=32,ascii"`
	Password string  `validate:"min=8,max=72"`
	UserInfo UserInfo
}

type Profile struct {
	Id        uuid.UUID
	Nickname  *string
	Email     string
	Gender    *string
	Birthday  *time.Time
	CreatedAt time.Time
}

func mapProfile(usr *models.User) Profile {
	return Profile{
		Id:        usr.Id,
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Gender:    usr.Gender,
		Birthday:  usr.Birthday,
		CreatedAt: usr.CreatedAt,
	}
}
