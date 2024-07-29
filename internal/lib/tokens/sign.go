package tokens

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration = 0
	RefreshTTL    time.Duration = 0
)

type Access struct {
	UserId    string `json:"uid"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

func (a Access) Valid() error {
	return nil
}

type Refresh struct {
	UserId    string `json:"uid"`
	Rotation  uint64 `json:"rtn"`
	ExpiresAt int64  `json:"exp"`
}

func (a Refresh) Valid() error {
	return nil
}

func (a *Access) Sign() (AccessString, error) {
	return a.SignSecret(AccessSecret)
}

func (a *Access) SignSecret(secret []byte) (AccessString, error) {
	if AccessTTL != 0 {
		a.ExpiresAt = time.Now().Add(AccessTTL).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a)

	tknstr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return AccessString(tknstr), nil
}

func (r *Refresh) Sign() (RefreshString, error) {
	return r.SignSecret(RefreshSecret)
}

func (r *Refresh) SignSecret(secret []byte) (RefreshString, error) {
	if AccessTTL != 0 {
		r.ExpiresAt = time.Now().Add(RefreshTTL).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, r)

	tknstr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return RefreshString(tknstr), nil
}
