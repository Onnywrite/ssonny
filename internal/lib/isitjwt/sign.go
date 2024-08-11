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

var TODOSecret = strings.Repeat("isitjwt?", 4)

type keys struct {
	pub  ed25519.PublicKey
	priv ed25519.PrivateKey
}

func Sign(secret string, userId uuid.UUID, subject string, exp time.Duration) (string, error) {
	body := userId.String() + ":" + subject + ":" + strconv.FormatInt(time.Now().Add(exp).Unix(), 10) + ":1"

	basedBody := base64.RawStdEncoding.EncodeToString([]byte(body))
	_, key, err := genKeys(secret)
	if err != nil {
		return "", err
	}
	sig := ed25519.Sign([]byte(key), []byte(basedBody))

	return basedBody + "." + base64.RawStdEncoding.EncodeToString(sig), nil
}

var (
	keysMap sync.Map
)

func genKeys(secret string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if len(secret) < 32 {
		return nil, nil, ErrSecretTooShort
	}
	if storedKeys, ok := keysMap.Load(secret); ok {
		pubPriv := storedKeys.(keys)
		return pubPriv.pub, pubPriv.priv, nil
	}

	secretReader := strings.NewReader(secret)
	pub, key, err := ed25519.GenerateKey(secretReader)
	if err != nil {
		return nil, nil, fmt.Errorf("error while generating keys: %w", err)
	}

	keysMap.Store(secret, keys{
		pub:  pub,
		priv: key,
	})

	return pub, key, nil
}
