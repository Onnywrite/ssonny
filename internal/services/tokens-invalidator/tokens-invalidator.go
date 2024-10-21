package tokensinvalidator

import (
	"context"
	"time"

	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/internal/domain/models"
)

type Service struct {
	frequency time.Duration
	repo      TokenRepo
}
type TokenRepo interface {
	Token(context.Context, uint64) (*models.Token, error)
}

type Dependencies struct {
	TokenRepo TokenRepo
}

type Config struct {
	Dependencies
	Frequency time.Duration
}

func New(deps Dependencies) *Service {
	conf := config.Get()

	return &Service{
		frequency: conf.Limits.TokensInvalidationFrequency,
		repo:      deps.TokenRepo,
	}
}

func NewWithConfig(c Config) *Service {
	return &Service{
		frequency: c.Frequency,
		repo:      c.TokenRepo,
	}
}
