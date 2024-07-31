package auth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/rs/zerolog"
)

type Service struct {
	log          *zerolog.Logger
	repo         UserRepo
	sessionRepo  SessionRepo
	emailService EmailService
}

type UserRepo interface {
	SaveUser(context.Context, models.User) (*models.User, repo.Transactor, error)
	UpdateUser(context.Context, models.User) error
	UserByEmail(context.Context, string) (*models.User, error)
	UserByNickname(context.Context, string) (*models.User, error)
}

type SessionRepo interface {
	SaveSession(context.Context, models.Session) (repo.Transactor, error)
}

type EmailService interface {
	SendVerificationEmail(context.Context, email.VerificationEmail) error
}

func NewService(log *zerolog.Logger, userRepo UserRepo, sessionRepo SessionRepo, emailService EmailService) *Service {
	return &Service{
		log:          log,
		repo:         userRepo,
		sessionRepo:  sessionRepo,
		emailService: emailService,
	}
}
