package auth

import (
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"

	"github.com/google/uuid"
)

type UserInfo struct {
	Platform string
	Agent    string
}

type RegisterWithPasswordData struct {
	Nickname *string
	Email    string
	Gender   *string
	Birthday *time.Time
	Password string
	UserInfo UserInfo
}

type AuthenticatedUser struct {
	Access  string
	Refresh string
	Profile Profile
}

type LoginWithPasswordData struct {
	Email    *string
	Nickname *string
	Password string
	UserInfo UserInfo
}

type Profile struct {
	Id        uuid.UUID
	Nickname  *string
	Email     string
	Gender    *string
	Verified  bool
	Birthday  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func mapProfile(usr *models.User) Profile {
	return Profile{
		Id:        usr.Id,
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Gender:    usr.Gender,
		Verified:  usr.Verified,
		Birthday:  usr.Birthday,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}
}
