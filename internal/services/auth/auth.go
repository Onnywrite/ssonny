package auth

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	log          zerolog.Logger
	repo         UserRepo
	emailService EmailService
	signer       TokenSigner
	tokenRepo    TokenRepo
}

type UserRepo interface {
	SaveUser(context.Context, models.User) (*models.User, repo.Transactor, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, newValues map[string]any) error
	UserByEmail(context.Context, string) (*models.User, error)
	UserByNickname(context.Context, string) (*models.User, error)
	UserById(context.Context, uuid.UUID) (*models.User, error)
}

type TokenRepo interface {
	SaveToken(context.Context, models.Token) (uint64, repo.Transactor, error)
	UpdateToken(ctx context.Context, id uint64, newValues map[string]any) error
	Token(context.Context, uint64) (*models.Token, error)
	DeleteTokens(context.Context, uuid.UUID, *uint64) error
	DeleteToken(context.Context, uint64) error
}

type TokenSigner interface {
	SignAccess(userId uuid.UUID, aud *uint64,
		authzParty string, scopes ...string) (string, error)
	SignRefresh(userId uuid.UUID, aud *uint64,
		authzParty string, rotation, jwtId uint64) (string, error)
	SignEmail(userId uuid.UUID) (string, error)
}

type EmailService interface {
	SendVerificationEmail(context.Context, email.VerificationEmail) error
}

type Config struct {
	UserRepo     UserRepo
	EmailService EmailService
	TokenRepo    TokenRepo
	TokensSigner TokenSigner
}

func NewService(log zerolog.Logger, c Config) *Service {
	return &Service{
		log:          log,
		repo:         c.UserRepo,
		emailService: c.EmailService,
		signer:       c.TokensSigner,
		tokenRepo:    c.TokenRepo,
	}
}
