package isitjwt

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// nolint: gochecknoglobals mnd
var TODOSecret = strings.Repeat("isitjwt?", 4)

type keys struct {
	pub  ed25519.PublicKey
	priv ed25519.PrivateKey
}

func Sign(secret string, userId uuid.UUID, subject string, exp time.Duration) (string, error) {
	// nolint: lll
	body := userId.String() + ":" + subject + ":" + strconv.FormatInt(time.Now().Add(exp).Unix(), 10) + ":1"
	basedBody := base64.RawStdEncoding.EncodeToString([]byte(body))

	_, key, err := genKeys(secret)
	if err != nil {
		return "", err
	}

	sig := ed25519.Sign([]byte(key), []byte(basedBody))

	return basedBody + "." + base64.RawStdEncoding.EncodeToString(sig), nil
}

// nolint: gochecknoglobals
var keysMap sync.Map

func genKeys(secret string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	const ed25519MinimumSecretLength = 32

	if len(secret) < ed25519MinimumSecretLength {
		return nil, nil, ErrSecretTooShort
	}

	if storedKeys, ok := keysMap.Load(secret); ok {
		pubPriv, ok := storedKeys.(keys)
		if !ok {
			return nil, nil, fmt.Errorf("error while retrieving keys: value is not of type keys")
		}

		return pubPriv.pub, pubPriv.priv, nil
	}

	pub, key, err := ed25519.GenerateKey(strings.NewReader(secret))
	if err != nil {
		return nil, nil, fmt.Errorf("error while generating keys: %w", err)
	}

	keysMap.Store(secret, keys{
		pub:  pub,
		priv: key,
	})

	return pub, key, nil
}
