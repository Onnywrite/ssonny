package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSecretTooShort      = errors.New("secret is too short, must be at least 32 bytes")
	ErrInvalidTokenVersion = errors.New("invalid version, expected 1")
	ErrInvalidSubject      = errors.New("invalid subject")
)

func (g Generator) SignEmail(userId uuid.UUID) (string, error) {
	return g.signMyToken(userId.String(), "email", g.emailExp, 1)
}

func (g Generator) ParseEmail(token string) (uuid.UUID, error) {
	id, err := g.parseMyToken(g.secret, "email", token)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}

func (g Generator) parseMyToken(secret []byte, subject, token string) (string, error) {
	const tokenPartsCount = 2

	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != tokenPartsCount {
		return "", ErrInvalidToken
	}

	body, rightSignature := splittedToken[0], splittedToken[1]

	debasedRightSignature, err := base64.RawURLEncoding.DecodeString(rightSignature)
	if err != nil {
		return "", fmt.Errorf("error while decoding signature: %w: %w", ErrInvalidToken, err)
	}

	signature := hmac.New(sha256.New, secret)
	signature.Write([]byte(body))

	if !hmac.Equal(debasedRightSignature, signature.Sum(nil)) {
		return "", ErrInvalidToken
	}

	debasedBody, err := base64.RawURLEncoding.DecodeString(body)
	if err != nil {
		return "", fmt.Errorf("error while decoding body: %w: %w", ErrInvalidToken, err)
	}

	splittedBody := strings.Split(string(debasedBody), ":")
	if splittedBody[len(splittedBody)-1] != "1" {
		return "", fmt.Errorf("%w: %w", ErrInvalidToken, ErrInvalidTokenVersion)
	}

	payload, sub, expStr := splittedBody[0], splittedBody[1], splittedBody[2]
	exp, err := strconv.ParseInt(expStr, 10, 64)

	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if time.Now().After(time.Unix(exp, 0)) {
		return "", ErrExpired
	}

	if sub != subject {
		return "", fmt.Errorf("%w: %w", ErrInvalidToken, ErrInvalidSubject)
	}

	return payload, nil
}

func (g Generator) signMyToken(
	payload, subject string,
	exp time.Duration,
	version int,
) (string, error) {
	body := fmt.Sprintf("%s:%s:%d:%d", payload, subject, time.Now().Add(exp).Unix(), version)
	basedBody := base64.RawURLEncoding.EncodeToString([]byte(body))

	signature := hmac.New(sha256.New, []byte(g.secret))
	signature.Write([]byte(basedBody))
	sig := signature.Sum(nil)

	return basedBody + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}
