package isitjwt

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// nolint: cyclop
func Verify(secret, subject, token string) (uuid.UUID, error) {
	const tokenPartsCount = 2

	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != tokenPartsCount {
		return uuid.Nil, ErrInvalidToken
	}

	body, sig := splittedToken[0], splittedToken[1]
	pub, _, err := genKeys(secret)
	if err != nil {
		return uuid.Nil, err
	}

	debasedSig, err := base64.RawStdEncoding.DecodeString(sig)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while decoding signature: %w: %w", ErrInvalidToken, err)
	}

	if !ed25519.Verify(pub, []byte(body), debasedSig) {
		return uuid.Nil, ErrInvalidToken
	}

	debasedBody, err := base64.RawStdEncoding.DecodeString(body)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while decoding body: %w: %w", ErrInvalidToken, err)
	}

	splittedBody := strings.Split(string(debasedBody), ":")
	if splittedBody[len(splittedBody)-1] != "1" {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrInvalidToken, ErrInvalidTokenVersion)
	}

	id, sub, expStr := splittedBody[0], splittedBody[1], splittedBody[2]

	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if time.Now().After(time.Unix(exp, 0)) {
		return uuid.Nil, ErrTokenExpired
	}

	if sub != subject {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrInvalidToken, ErrInvalidSubject)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while parsing uuid: %w: %w", ErrInvalidToken, err)
	}

	return uid, nil
}
