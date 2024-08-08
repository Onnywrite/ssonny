package tokens

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/golang-jwt/jwt"
)

type Generator struct {
	issuer     string
	accessExp  time.Duration
	refreshExp time.Duration
	idExp      time.Duration
	pub        *rsa.PublicKey
	priv       *rsa.PrivateKey
	parser     jwt.Parser
}

func New(cfg config.TokensConfig) (Generator, error) {
	certPEM, err := tls.LoadX509KeyPair(cfg.PublicPath, cfg.SecretPath)
	if err != nil {
		return Generator{}, err
	}

	cert, err := x509.ParseCertificate(certPEM.Certificate[0])
	if err != nil {
		return Generator{}, err
	}

	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return Generator{}, err
	}

	privateKey, ok := certPEM.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return Generator{}, err
	}

	return NewWithKeys(cfg.Issuer, cfg.AccessTTL,
		cfg.RefreshTTL, cfg.IdTTL,
		publicKey, privateKey), nil
}

func NewWithKeys(iss string, aexp, rexp, iexp time.Duration,
	pub *rsa.PublicKey, priv *rsa.PrivateKey) Generator {
	return Generator{
		issuer:     iss,
		accessExp:  aexp,
		refreshExp: rexp,
		idExp:      iexp,
		pub:        pub,
		priv:       priv,
		parser: jwt.Parser{
			ValidMethods:         []string{"RS256"},
			UseJSONNumber:        true,
			SkipClaimsValidation: true,
		},
	}
}
