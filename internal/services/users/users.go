package users

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/services/email"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	log          zerolog.Logger
	repo         UserRepo
	emailService EmailService
}

type UserRepo interface {
	UpdateUser(ctx context.Context, userId uuid.UUID, newValues map[string]any) error
	UpdateAndGetUser(ctx context.Context, userId uuid.UUID, newValues map[string]any) (*models.User, error)
	UserById(context.Context, uuid.UUID) (*models.User, error)
}

type EmailService interface {
	SendPasswordResetEmail(context.Context, email.PasswordResetEmail) error
}

type Config struct {
	UserRepo     UserRepo
	EmailService EmailService
}

func NewService(log zerolog.Logger, c Config) *Service {
	return &Service{
		log:          log,
		repo:         c.UserRepo,
		emailService: c.EmailService,
	}
}
