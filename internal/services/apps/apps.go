package apps

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	log       *zerolog.Logger
	repo      AppRepo
	auth      AuthService
	validate  *validator.Validate
	generator RandomGenerator
}

type AppRepo interface {
	SaveApp(context.Context, models.App, []uint64) (*models.App, error)
	SaveDomains(context.Context, uuid.UUID, []models.Domain) ([]models.Domain, error)
	TieDomainsToApp(context.Context, uint64, []uint64) error
}

type AuthService interface {
	IsVerified(ctx context.Context, userId uuid.UUID) (bool, error)
}

type RandomGenerator interface {
	RandomString() (string, error)
}

func NewService(
	logger *zerolog.Logger,
	appRepo AppRepo,
	authService AuthService,
	randomGenerator RandomGenerator,
) *Service {
	return &Service{
		log:       logger,
		repo:      appRepo,
		auth:      authService,
		validate:  validator.New(validator.WithRequiredStructEnabled()),
		generator: randomGenerator,
	}
}
