package handlersusers

import (
	"context"
	"time"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/services/users"

	"github.com/google/uuid"
)

type UsersService interface {
	GetProfile(ctx context.Context, userId uuid.UUID) (*users.Profile, error)
	PutProfile(ctx context.Context, userId uuid.UUID,
		data users.UpdateProfileData) (*users.Profile, error)
	PutProfilePassword(ctx context.Context, userId uuid.UUID,
		data users.UpdatePasswordData) error
}

type UsersHandler struct {
	Service UsersService
}

func (h *UsersHandler) GetProfile(ctx context.Context,
	request api.GetProfileRequestObject,
) (api.GetProfileResponseObject, error) {
	userId := ctx.Value("currentUserId").(uuid.UUID)

	profile, err := h.Service.GetProfile(ctx, userId)
	if err != nil {
		return getGetProfileResponse(err)
	}

	return api.GetProfile200JSONResponse(toApiProfile(*profile)), nil
}

func getGetProfileResponse(serviceError error,
) (api.GetProfileResponseObject, error) {
	if erix.HttpCode(serviceError) == erix.CodeNotFound {
		return api.GetProfile404JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceError.Error(),
		}, nil
	}

	return nil, serviceError
}

func (h *UsersHandler) PutProfile(ctx context.Context,
	request api.PutProfileRequestObject,
) (api.PutProfileResponseObject, error) {
	userId := ctx.Value("currentUserId").(uuid.UUID)

	profile, err := h.Service.PutProfile(ctx, userId, users.UpdateProfileData{
		Birthday: request.Body.Birthday,
		Gender:   request.Body.Gender,
		Nickname: request.Body.Nickname,
	})
	if err != nil {
		return getPutProfileResponse(err)
	}

	return api.PutProfile200JSONResponse(toApiProfile(*profile)), nil
}

func getPutProfileResponse(serviceErr error,
) (api.PutProfileResponseObject, error) {
	if erix.HttpCode(serviceErr) == erix.CodeNotFound {
		return api.PutProfile404JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceErr.Error(),
		}, nil
	}

	if erix.HttpCode(serviceErr) == erix.CodeConflict {
		return api.PutProfile409JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceErr.Error(),
		}, nil
	}

	return nil, serviceErr
}

func toApiProfile(profile users.Profile) api.Profile {
	var birthdayString *string

	if profile.Birthday != nil {
		birthdayStringLocal := profile.Birthday.Format(time.DateOnly)
		birthdayString = &birthdayStringLocal
	}

	return api.Profile{
		Id:        profile.Id,
		Nickname:  profile.Nickname,
		Email:     profile.Email,
		Gender:    profile.Gender,
		Verified:  profile.Verified,
		Birthday:  birthdayString,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}
}
