package handlersauth

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
)

func (h *AuthHandler) GetAuthCheck(ctx context.Context,
	request api.GetAuthCheckRequestObject,
) (api.GetAuthCheckResponseObject, error) {
	return api.GetAuthCheck200Response{}, nil
}
