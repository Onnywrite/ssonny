package handlersauth

import (
	"context"
	"strings"
	"time"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/lib/erix"
	"github.com/Onnywrite/ssonny/internal/services/auth"

	"github.com/mileusna/useragent"
)

type AuthService interface {
	Registrator
	Loginer
	Logouter
	Refresher
	Verifier
}

type Registrator interface {
	RegisterWithPassword(ctx context.Context,
		data auth.RegisterWithPasswordData) (*auth.AuthenticatedUser, error)
}

// I will use generated code in the service layer, but
// now I decided te leave everything as it is just to
// KEEP IT SIMPLE :-)
//
// Otherwise I'd have to reweite the whole service layer INCLUDE tests,
// which does not make me happy.
type AuthHandler struct {
	Service          AuthService
	RefreshParser    RefreshTokenParser
	EmailTokenParser EmailTokenParser
}

func (h *AuthHandler) PostAuthRegisterWithPassword(ctx context.Context,
	request api.PostAuthRegisterWithPasswordRequestObject,
) (api.PostAuthRegisterWithPasswordResponseObject, error) {
	var birthday *time.Time

	if request.Body.Birthday != nil {
		bday, err := time.Parse(time.DateOnly, *request.Body.Birthday)
		if err != nil {
			return api.PostAuthRegisterWithPassword400JSONResponse{ //nolint: nilerr
				Service: api.ValidationErrorServiceSsonny,
				Fields:  map[string]any{"birthday": "Birthday has invalid date format"},
			}, nil
		}

		birthday = &bday
	}

	authUser, err := h.Service.RegisterWithPassword(ctx, auth.RegisterWithPasswordData{
		Nickname: request.Body.Nickname,
		Email:    request.Body.Email,
		Gender:   request.Body.Gender,
		Birthday: birthday,
		Password: request.Body.Password,
		UserInfo: getUserInfo(request.Params.UserAgent),
	})
	if err != nil {
		return getPostAuthRegisterWithPasswordResponse(err)
	}

	return api.PostAuthRegisterWithPassword201JSONResponse{
		Access:  authUser.Access,
		Refresh: authUser.Refresh,
		Profile: toApiProfile(authUser.Profile),
	}, nil
}

func getPostAuthRegisterWithPasswordResponse(serviceError error,
) (api.PostAuthRegisterWithPasswordResponseObject, error) {
	if erix.HttpCode(serviceError) == erix.CodeConflict {
		return api.PostAuthRegisterWithPassword409JSONResponse{
			Service: api.ErrServiceSsonny,
			Message: serviceError.Error(),
		}, nil
	}

	return nil, serviceError
}

func getUserInfo(userAgent string) auth.UserInfo {
	ua := useragent.Parse(userAgent)
	platform := strings.Join([]string{ua.OS, ua.OSVersion}, " ")
	agent := strings.Join([]string{ua.Name, ua.Version}, " ")

	return auth.UserInfo{
		Platform: platform,
		Agent:    agent,
	}
}

func toApiProfile(profile auth.Profile) api.Profile {
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
	}
}
