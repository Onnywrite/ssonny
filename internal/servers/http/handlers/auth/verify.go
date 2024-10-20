package handlersauth

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
	userId, err := h.EmailTokenParser.ParseEmail(request.Body.Token)
	if err != nil {
		return api.PostAuthVerifyEmail400JSONResponse{ //nolint: nilerr
			Service: api.ValidationErrorServiceSsonny,
			Fields: map[string]any{
				"Token": err.Error(),
			},
		}, nil
	}

	if err := h.Service.VerifyEmail(ctx, userId); err != nil {
		return nil, err
	}

	return api.PostAuthVerifyEmail200Response{}, nil
}
