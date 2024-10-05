package httpserver

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
	"github.com/Onnywrite/ssonny/internal/servers/http/handlers"

	// authh "github.com/Onnywrite/ssonny/internal/servers/http/handlers/auth"

	"github.com/gofiber/fiber/v3"
)

// type AuthService interface {
// 	authh.Registrator
// 	authh.Loginer
// 	authh.Logouter
// 	authh.Refresher
// 	authh.Verifier
// }

// type TokenParser interface {
// 	middlewares.AccessTokenParser
// 	authh.RefreshTokenParser
// 	authh.EmailTokenParser
// }

func InitApi(r fiber.Router) { //authService AuthService, tokenParser TokenParser) {
	r.Get("/ping", handlers.Ping())

	sr := api.NewStrictHandler(&handler{}, nil)
	api.RegisterHandlersWithOptions(r, sr, api.FiberServerOptions{
		BaseURL:     "/api",
		Middlewares: []api.MiddlewareFunc{},
	})
}

type handler struct{}

func (h *handler) PostAuthLoginWithPassword(ctx context.Context,
	request api.PostAuthLoginWithPasswordRequestObject,
) (api.PostAuthLoginWithPasswordResponseObject, error) {
	return api.PostAuthLoginWithPassword200JSONResponse{}, nil
}

func (h *handler) PostAuthLogout(ctx context.Context,
	request api.PostAuthLogoutRequestObject,
) (api.PostAuthLogoutResponseObject, error) {
	return api.PostAuthLogout200Response{}, nil
}

func (h *handler) PostAuthRefresh(ctx context.Context,
	request api.PostAuthRefreshRequestObject,
) (api.PostAuthRefreshResponseObject, error) {
	return api.PostAuthRefresh200JSONResponse{}, nil
}

func (h *handler) PostAuthRegisterWithPassword(ctx context.Context,
	request api.PostAuthRegisterWithPasswordRequestObject,
) (api.PostAuthRegisterWithPasswordResponseObject, error) {
	return api.PostAuthRegisterWithPassword201JSONResponse{}, nil
}

func (h *handler) PostAuthVerifyEmail(ctx context.Context,
	request api.PostAuthVerifyEmailRequestObject,
) (api.PostAuthVerifyEmailResponseObject, error) {
	return api.PostAuthVerifyEmail200Response{}, nil
}

func (h *handler) GetHealthz(ctx context.Context,
	request api.GetHealthzRequestObject,
) (api.GetHealthzResponseObject, error) {
	return api.GetHealthz200TextResponse("ok"), nil
}

func (h *handler) GetMetrics(ctx context.Context,
	request api.GetMetricsRequestObject,
) (api.GetMetricsResponseObject, error) {
	return api.GetMetrics200JSONResponse{}, nil
}

func (h *handler) GetPing(ctx context.Context,
	request api.GetPingRequestObject,
) (api.GetPingResponseObject, error) {
	return api.GetPing200TextResponse("pong"), nil
}
