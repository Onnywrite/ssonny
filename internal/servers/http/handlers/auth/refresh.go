package handlersapiauth

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/lib/tokens"
	"github.com/Onnywrite/ssonny/internal/services/auth"
)

type Refresher interface {
	Refresh(ctx context.Context, token tokens.Refresh) (*auth.Tokens, error)
}

func (h *AuthHandler) PostAuthRefresh(ctx context.Context,
	request api.PostAuthRefreshRequestObject,
) (api.PostAuthRefreshResponseObject, error) {
	parserRefresh, err := h.RefreshParser.ParseRefresh(request.Body.RefreshToken)
	if err != nil {
		return api.PostAuthRefresh401JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: err.Error(),
		}, nil
	}

	tokens, err := h.Service.Refresh(ctx, *parserRefresh)
	if err != nil {
		return getPostAuthRefreshResponse(err)
	}

	return api.PostAuthRefresh200JSONResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func getPostAuthRefreshResponse(serviceError error,
) (api.PostAuthRefreshResponseObject, error) {
	if erix.HttpCode(serviceError) == erix.CodeUnauthorized {
		return api.PostAuthRefresh401JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceError.Error(),
		}, nil
	}

	return nil, serviceError
}
