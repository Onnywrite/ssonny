package auth

import (
	"context"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	Platform string
	Agent    string
}

type RegisterWithPasswordData struct {
	Nickname string
	Email    string
	Gender   string
	Birthday *time.Time
	Password string
	UserInfo UserInfo
}

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

type AuthenticatedUser struct {
	tokens.Pair
	Profile Profile
}

// RegisterWithPassword registrates new user with unique email and unique nickname
func (s *Service) RegisterWithPassword(ctx context.Context, data RegisterWithPasswordData) (*AuthenticatedUser, error) {
	log := s.log.With().Str("user_nickname", data.Nickname).Str("user_email", data.Email).Logger()
	// TODO: validate data

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("error while hashing password")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	saved, err := s.repo.SaveUser(ctx, models.User{
		Nickname:     data.Nickname,
		Email:        data.Email,
		Gender:       data.Gender,
		Birthday:     data.Birthday,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, userFailed(&log, err)
	}

	return s.beginSession(ctx, saved, &data.UserInfo)
}

func (s *Service) beginSession(ctx context.Context, saved *models.User, info *UserInfo) (*AuthenticatedUser, error) {
	tx, err := s.sessionRepo.SaveSession(ctx, models.Session{
		UserId:       saved.Id,
		AppId:        0,
		LastRotation: 0,
		Platform:     info.Platform,
		Agent:        info.Agent,
	})
	if err != nil {
		return nil, sessionFailed(s.log, err)
	}
	defer tx.Rollback()

	pair, err := tokens.NewPair(saved, 0)
	if err != nil {
		s.log.Error().Err(err).Msg("error while creating tokens")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error while committing session saving")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return &AuthenticatedUser{
		Pair:    pair,
		Profile: mapProfile(saved),
	}, nil
}
