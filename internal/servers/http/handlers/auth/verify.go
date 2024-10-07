package handlersapiauth

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/google/uuid"
)

type Verifier interface {
	VerifyEmail(ctx context.Context, userId uuid.UUID) error
}

type EmailTokenParser interface {
	ParseEmail(token string) (uuid.UUID, error)
}

func (h *AuthHandler) PostAuthVerifyEmail(ctx context.Context,
	request api.PostAuthVerifyEmailRequestObject,
) (api.PostAuthVerifyEmailResponseObject, error) {
	userId, err := h.EmailTokenParser.ParseEmail(request.Params.Token)
	if err != nil {
		return api.PostAuthVerifyEmail400JSONResponse{
			Service: api.ValidationErrorServiceSsonny,
			Fields: map[string]any{
				"token": "invalid token",
			},
		}, nil
	}

	if err := h.Service.VerifyEmail(ctx, userId); err != nil {
		return nil, err
	}

	return api.PostAuthVerifyEmail200Response{}, nil
}
