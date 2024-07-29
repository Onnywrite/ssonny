package auth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/rs/zerolog"
)

type Service struct {
	log         *zerolog.Logger
	repo        UserRepo
	sessionRepo SessionRepo
}

type UserRepo interface {
	SaveUser(context.Context, models.User) (*models.User, error)
}

type SessionRepo interface {
	SaveSession(context.Context, models.Session) (repo.Transactor, error)
}

func NewService(log *zerolog.Logger, userRepo UserRepo, sessionRepo SessionRepo) *Service {
	return &Service{
		log:         log,
		repo:        userRepo,
		sessionRepo: sessionRepo,
	}
}
