package handlersusers

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/services/users"
	"github.com/google/uuid"
)

func (h *UsersHandler) PutProfilePassword(ctx context.Context,
	request api.PutProfilePasswordRequestObject,
) (api.PutProfilePasswordResponseObject, error) {
	userId := ctx.Value("currentUserId").(uuid.UUID)

	err := h.Service.PutProfilePassword(ctx, userId, users.UpdatePasswordData{
		CurrentPassword: request.Body.CurrentPassword,
		NewPassword:     request.Body.NewPassword,
	})
	if err != nil {
		return getPutProfilePasswordResponse(err)
	}

	return api.PutProfilePassword200Response{}, nil
}

func getPutProfilePasswordResponse(serviceErr error,
) (api.PutProfilePasswordResponseObject, error) {
	if erix.HttpCode(serviceErr) == erix.CodeNotFound {
		return api.PutProfilePassword404JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceErr.Error(),
		}, nil
	}

	return nil, serviceErr
}
