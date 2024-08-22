package tokens

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func (g Generator) ParseAccess(token string) (*Access, error) {
	var claims Access

	_, err := g.parser.ParseWithClaims(token, &claims, func(*jwt.Token) (interface{}, error) {
		return g.pub, nil
	})
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrExpired
	}

	if claims.Issuer != g.issuer {
		return nil, ErrIssuerMismatch
	}

	return &claims, nil
}

func (g Generator) ParseRefresh(token string) (*Refresh, error) {
	var claims Refresh

	_, err := g.parser.ParseWithClaims(token, &claims, func(*jwt.Token) (interface{}, error) {
		return g.pub, nil
	})
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrExpired
	}

	if claims.Issuer != g.issuer {
		return nil, ErrIssuerMismatch
	}

	return &claims, nil
}

func (g Generator) ParseId(token string) (*Id, error) {
	var claims Id

	_, err := g.parser.ParseWithClaims(token, &claims, func(*jwt.Token) (interface{}, error) {
		return g.pub, nil
	})
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrExpired
	}

	if claims.Issuer != g.issuer {
		return nil, ErrIssuerMismatch
	}

	return &claims, nil
}
