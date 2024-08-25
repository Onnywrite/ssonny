package tokens

import (
	"time"

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

func New(iss, secret string, accessTtl, refreshTtl, idExp, emailEcp time.Duration) Generator {
	return Generator{
		issuer:     iss,
		accessExp:  accessTtl,
		refreshExp: refreshTtl,
		idExp:      idExp,
		emailExp:   emailEcp,
		secret:     []byte(secret),
		parser: jwt.Parser{
			ValidMethods:         []string{"HS256"},
			UseJSONNumber:        true,
			SkipClaimsValidation: true,
		},
	}
}
