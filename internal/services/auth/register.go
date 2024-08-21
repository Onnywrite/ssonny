package auth

import (
	"context"
	"strings"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/isitjwt"
	"github.com/Onnywrite/ssonny/internal/services/email"
	"golang.org/x/crypto/bcrypt"
)

// RegisterWithPassword registrates new user with unique email and unique nickname
func (s *Service) RegisterWithPassword(ctx context.Context, data RegisterWithPasswordData) (*AuthenticatedUser, error) {
	log := s.log.With().Str("user_email", data.Email).Logger()
	if err := data.Validate(); err != nil {
		log.Debug().Err(err).Msg("invalid data, bad request")
		return nil, erix.Wrap(err, erix.CodeBadRequest, ErrInvalidData)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("error while hashing password")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	stringHash := string(hash)
	saved, tx, err := s.repo.SaveUser(ctx, models.User{
		Nickname:     data.Nickname,
		Email:        data.Email,
		IsVerified:   false,
		Gender:       data.Gender,
		Birthday:     data.Birthday,
		PasswordHash: &stringHash,
	})
	if err != nil {
		return nil, userFailed(&log, err)
	}
	// nolint: errcheck
	defer tx.Rollback()

	// TODO: configs for this
	token, err := isitjwt.Sign(isitjwt.TODOSecret, saved.Id, SubjectEmail, time.Hour*2)
	if err != nil {
		log.Error().Err(err).Msg("error while signing email verification token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	var userNickname string
	if data.Nickname != nil {
		userNickname = *data.Nickname
	} else {
		userNickname = strings.Split(data.Email, "@")[0]
	}
	err = s.emailService.SendVerificationEmail(ctx, email.VerificationEmail{
		Recipient:    data.Email,
		UserNickname: userNickname,
		Token:        token,
	})
	if err != nil {
		log.Warn().Err(err).Msg("email has not been sent")
		return nil, erix.Wrap(err, erix.CodeBadRequest, ErrEmailUnverified)
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error while committing session saving")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return s.generateAndSaveTokens(ctx, *saved, data.UserInfo)
}

func (s *Service) generateAndSaveTokens(ctx context.Context, user models.User, info UserInfo) (*AuthenticatedUser, error) {
	log := s.log.With().Str("user_email", user.Email).Logger()

	access, err := s.signer.SignAccess(user.Id, nil, "self", "*")
	if err != nil {
		log.Error().Err(err).Msg("error while signing access token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	jwtId, tx, err := s.tokenRepo.SaveToken(ctx, models.Token{
		UserId:    user.Id,
		AppId:     nil,
		Rotation:  0,
		RotatedAt: time.Now(),
		Platform:  info.Platform,
		Agent:     info.Agent,
	})
	if err != nil {
		log.Error().Err(err).Msg("error while saving token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}
	// nolint: errcheck
	defer tx.Rollback()

	refresh, err := s.signer.SignRefresh(user.Id, nil, "self", 0, jwtId)
	if err != nil {
		log.Error().Err(err).Msg("error while signing refresh token")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error while committing session saving")
		return nil, erix.Wrap(err, erix.CodeInternalServerError, ErrInternal)
	}

	return &AuthenticatedUser{
		Access:  access,
		Refresh: refresh,
		Profile: mapProfile(&user),
	}, nil
}
