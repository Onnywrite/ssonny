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
	parser     jwt.Parser
	secret     string
}

func New(iss, secret string, accessTtl, refreshTtl, idTxp time.Duration) Generator {
	return Generator{
		issuer:     iss,
		accessExp:  accessTtl,
		refreshExp: refreshTtl,
		idExp:      idTxp,
		secret:     secret,
		parser: jwt.Parser{
			ValidMethods:         []string{"HS256"},
			UseJSONNumber:        true,
			SkipClaimsValidation: true,
		},
	}
}
