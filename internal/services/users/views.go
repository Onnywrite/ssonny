package users

import (
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"

	"github.com/google/uuid"
)

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

type UpdateProfileData struct {
	Birthday *string
	Gender   *string
	Nickname *string
}
