package auth

import (
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/google/uuid"
)

type Profile struct {
	Id        uuid.UUID
	Nickname  string
	Email     string
	Gender    string // default, I guess, 'not specified'
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
