package handlersauth

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/services/auth"
)

type Loginer interface {
	LoginWithPassword(ctx context.Context,
		data auth.LoginWithPasswordData) (*auth.AuthenticatedUser, error)
}

func (h *AuthHandler) PostAuthLoginWithPassword(ctx context.Context,
	request api.PostAuthLoginWithPasswordRequestObject,
) (api.PostAuthLoginWithPasswordResponseObject, error) {
	userInfo := getUserInfo(request.Params.UserAgent)

	authUser, err := h.Service.LoginWithPassword(ctx, auth.LoginWithPasswordData{
		Email:    request.Body.Email,
		Nickname: request.Body.Nickname,
		Password: request.Body.Password,
		UserInfo: userInfo,
	})
	if err != nil {
		return getPostAuthLoginWithPasswordResponse(err)
	}

	return api.PostAuthLoginWithPassword200JSONResponse{
		Access:  authUser.Access,
		Refresh: authUser.Refresh,
		Profile: toApiProfile(authUser.Profile),
	}, nil
}

func getPostAuthLoginWithPasswordResponse(serviceError error,
) (api.PostAuthLoginWithPasswordResponseObject, error) {
	if erix.HttpCode(serviceError) == erix.CodeNotFound {
		return api.PostAuthLoginWithPassword404JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceError.Error(),
		}, nil
	}

	return nil, serviceError
}
