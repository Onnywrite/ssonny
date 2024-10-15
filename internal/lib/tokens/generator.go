package tokens

import (
	"time"

	"github.com/Onnywrite/ssonny/internal/config"

	"github.com/golang-jwt/jwt"
)

type Generator struct {
	issuer     string
	accessExp  time.Duration
	refreshExp time.Duration
	idExp      time.Duration
	emailExp   time.Duration
	parser     jwt.Parser
	secret     []byte
}

type Config struct {
	Issuer     string
	AccessExp  time.Duration
	RefreshExp time.Duration
	IdExp      time.Duration
	EmailExp   time.Duration
	Secret     []byte
}

func New() Generator {
	conf := config.Get()

	return NewWithConfig(Config{
		Issuer:     conf.Tokens.Issuer,
		AccessExp:  conf.Tokens.AccessTtl,
		RefreshExp: conf.Tokens.RefreshTtl,
		IdExp:      conf.Tokens.IdTtl,
		EmailExp:   conf.Tokens.EmailVerificationTtl,
		Secret:     []byte(conf.Secrets.SecretString),
	})
}

func NewWithConfig(c Config) Generator {
	return Generator{
		issuer:     c.Issuer,
		accessExp:  c.AccessExp,
		refreshExp: c.RefreshExp,
		idExp:      c.IdExp,
		emailExp:   c.EmailExp,
		secret:     c.Secret,
		parser: jwt.Parser{
			ValidMethods:         []string{"HS256"},
			UseJSONNumber:        true,
			SkipClaimsValidation: true,
		},
	}
}
