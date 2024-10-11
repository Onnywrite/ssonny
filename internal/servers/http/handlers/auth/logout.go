package handlersauth

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
)

type Logouter interface {
	Logout(ctx context.Context, jwtId uint64) error
}

type RefreshTokenParser interface {
	ParseRefresh(token string) (*tokens.Refresh, error)
}

func (h *AuthHandler) PostAuthLogout(ctx context.Context,
	request api.PostAuthLogoutRequestObject,
) (api.PostAuthLogoutResponseObject, error) {
	parsedRefresh, err := h.RefreshParser.ParseRefresh(request.Body.RefreshToken)
	if err != nil {
		return api.PostAuthLogout401JSONResponse{ //nolint: nilerr
			Service: api.ErrServiceSsonny,
			Message: err.Error(),
		}, nil
	}

	if err := h.Service.Logout(ctx, parsedRefresh.Id); err != nil {
		// there is only 500 status
		return nil, err
	}

	return api.PostAuthLogout200Response{}, nil
}
