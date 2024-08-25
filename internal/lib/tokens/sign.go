package tokens

import (
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (g Generator) SignAccess(userId uuid.UUID,
	aud *uint64,
	authzParty string,
	scopes ...string,
) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, Access{
		Issuer:          g.issuer,
		Subject:         userId,
		Audience:        aud,
		AuthorizedParty: authzParty,
		ExpiresAt:       time.Now().Add(g.accessExp).Unix(),
		Scopes:          scopes,
	})

	token, err := tkn.SignedString(g.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (g Generator) SignRefresh(userId uuid.UUID,
	aud *uint64,
	authzParty string,
	rotation, jwtId uint64,
) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, Refresh{
		Issuer:          g.issuer,
		Subject:         userId,
		Audience:        aud,
		AuthorizedParty: authzParty,
		ExpiresAt:       time.Now().Add(g.refreshExp).Unix(),
		Id:              jwtId,
		Rotation:        rotation,
	})

	token, err := tkn.SignedString(g.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (g Generator) SignId(user *models.User,
	aud *uint64,
	authzParty string,
	roles ...string,
) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, Id{
		Issuer:          g.issuer,
		Subject:         user.Id,
		Audience:        aud,
		AuthorizedParty: authzParty,
		ExpiresAt:       time.Now().Add(g.idExp).Unix(),
		Nickname:        user.Nickname,
		Email:           user.Email,
		IsVerified:      user.IsVerified,
		Gender:          user.Gender,
		Birthday:        user.Birthday,
		Roles:           roles,
	})

	token, err := tkn.SignedString(g.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
