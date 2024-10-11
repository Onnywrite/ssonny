// Package httpapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package httpapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	fmtvalidate "github.com/Onnywrite/ssonny/pkg/fmtvalidate"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v3"
	"github.com/oapi-codegen/runtime"
) // ServerInterface represents all server handlers.
type ServerInterface interface {
	// Checks if access token is valid
	// (GET /auth/check)
	GetAuthCheck(c fiber.Ctx) error
	// Login user by their password and email or nickname
	// (POST /auth/loginWithPassword)
	PostAuthLoginWithPassword(c fiber.Ctx, params PostAuthLoginWithPasswordParams) error
	// Logouts user by invalidating refresh token
	// (POST /auth/logout)
	PostAuthLogout(c fiber.Ctx) error
	// Refreshes expired access and unexpired refresh tokens
	// (POST /auth/refresh)
	PostAuthRefresh(c fiber.Ctx) error
	// Registrates user by a password and email or nickname
	// (POST /auth/registerWithPassword)
	PostAuthRegisterWithPassword(c fiber.Ctx, params PostAuthRegisterWithPasswordParams) error
	// Verifies the user's email.
	// (POST /auth/verify/email)
	PostAuthVerifyEmail(c fiber.Ctx, params PostAuthVerifyEmailParams) error
	// The server's health probes
	// (GET /healthz)
	GetHealthz(c fiber.Ctx) error
	// OpenTelemetry metrics
	// (GET /metrics)
	GetMetrics(c fiber.Ctx) error
	// Pings the server
	// (GET /ping)
	GetPing(c fiber.Ctx) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// GetAuthCheck operation middleware
func (siw *ServerInterfaceWrapper) GetAuthCheck(c fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	return siw.Handler.GetAuthCheck(c)
}

// PostAuthLoginWithPassword operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLoginWithPassword(c fiber.Ctx) error {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthLoginWithPasswordParams

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "User-Agent" -------------
	if values, found := headers[http.CanonicalHeaderKey("User-Agent")]; found {
		var UserAgent UserAgent

		err = runtime.BindStyledParameterWithOptions("simple", "User-Agent", values[0], &UserAgent, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid format for parameter User-Agent: %w", err).Error())
		}

		params.UserAgent = UserAgent

	} else {
		err = fmt.Errorf("header parameter User-Agent is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.PostAuthLoginWithPassword(c, params)
}

// PostAuthLogout operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLogout(c fiber.Ctx) error {

	return siw.Handler.PostAuthLogout(c)
}

// PostAuthRefresh operation middleware
func (siw *ServerInterfaceWrapper) PostAuthRefresh(c fiber.Ctx) error {

	return siw.Handler.PostAuthRefresh(c)
}

// PostAuthRegisterWithPassword operation middleware
func (siw *ServerInterfaceWrapper) PostAuthRegisterWithPassword(c fiber.Ctx) error {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthRegisterWithPasswordParams

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "User-Agent" -------------
	if values, found := headers[http.CanonicalHeaderKey("User-Agent")]; found {
		var UserAgent UserAgent

		err = runtime.BindStyledParameterWithOptions("simple", "User-Agent", values[0], &UserAgent, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid format for parameter User-Agent: %w", err).Error())
		}

		params.UserAgent = UserAgent

	} else {
		err = fmt.Errorf("header parameter User-Agent is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.PostAuthRegisterWithPassword(c, params)
}

// PostAuthVerifyEmail operation middleware
func (siw *ServerInterfaceWrapper) PostAuthVerifyEmail(c fiber.Ctx) error {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthVerifyEmailParams

	var query url.Values
	query, err = url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "token" -------------

	if paramValue := c.Query("token"); paramValue != "" {

	} else {
		err = fmt.Errorf("query argument token is required, but not found")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = runtime.BindQueryParameter("form", true, true, "token", query, &params.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid format for parameter token: %w", err).Error())
	}

	return siw.Handler.PostAuthVerifyEmail(c, params)
}

// GetHealthz operation middleware
func (siw *ServerInterfaceWrapper) GetHealthz(c fiber.Ctx) error {

	return siw.Handler.GetHealthz(c)
}

// GetMetrics operation middleware
func (siw *ServerInterfaceWrapper) GetMetrics(c fiber.Ctx) error {

	return siw.Handler.GetMetrics(c)
}

// GetPing operation middleware
func (siw *ServerInterfaceWrapper) GetPing(c fiber.Ctx) error {

	return siw.Handler.GetPing(c)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	Middlewares         []fiber.Handler
	EndpointMiddlewares map[string][]fiber.Handler
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// created by github.com/Onnywrite
// Constants for all endpoints
const (
	// GET /auth/check: Checks if access token is valid
	EP_GetAuthCheck = "/auth/check"
	// POST /auth/loginWithPassword: Login user by their password and email or nickname
	EP_PostAuthLoginWithPassword = "/auth/loginWithPassword"
	// POST /auth/logout: Logouts user by invalidating refresh token
	EP_PostAuthLogout = "/auth/logout"
	// POST /auth/refresh: Refreshes expired access and unexpired refresh tokens
	EP_PostAuthRefresh = "/auth/refresh"
	// POST /auth/registerWithPassword: Registrates user by a password and email or nickname
	EP_PostAuthRegisterWithPassword = "/auth/registerWithPassword"
	// POST /auth/verify/email: Verifies the user's email.
	EP_PostAuthVerifyEmail = "/auth/verify/email"
	// GET /healthz: The server's health probes
	EP_GetHealthz = "/healthz"
	// GET /metrics: OpenTelemetry metrics
	EP_GetMetrics = "/metrics"
	// GET /ping: Pings the server
	EP_GetPing = "/ping"
)

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(m)
	}

	mws := func(ep string) []fiber.Handler {
		return options.EndpointMiddlewares[ep]
	}

	router.Get(EP_GetAuthCheck, wrapper.GetAuthCheck, mws(EP_GetAuthCheck)...)
	router.Post(EP_PostAuthLoginWithPassword, wrapper.PostAuthLoginWithPassword, mws(EP_PostAuthLoginWithPassword)...)
	router.Post(EP_PostAuthLogout, wrapper.PostAuthLogout, mws(EP_PostAuthLogout)...)
	router.Post(EP_PostAuthRefresh, wrapper.PostAuthRefresh, mws(EP_PostAuthRefresh)...)
	router.Post(EP_PostAuthRegisterWithPassword, wrapper.PostAuthRegisterWithPassword, mws(EP_PostAuthRegisterWithPassword)...)
	router.Post(EP_PostAuthVerifyEmail, wrapper.PostAuthVerifyEmail, mws(EP_PostAuthVerifyEmail)...)
	router.Get(EP_GetHealthz, wrapper.GetHealthz, mws(EP_GetHealthz)...)
	router.Get(EP_GetMetrics, wrapper.GetMetrics, mws(EP_GetMetrics)...)
	router.Get(EP_GetPing, wrapper.GetPing, mws(EP_GetPing)...)

}

type GetAuthCheckRequestObject struct {
}

type GetAuthCheckResponseObject interface {
	VisitGetAuthCheckResponse(ctx fiber.Ctx) error
}

type GetAuthCheck200Response struct {
}

func (response GetAuthCheck200Response) VisitGetAuthCheckResponse(ctx fiber.Ctx) error {
	ctx.Status(200)
	return nil
}

type GetAuthCheck401JSONResponse Err

func (response GetAuthCheck401JSONResponse) VisitGetAuthCheckResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type PostAuthLoginWithPasswordRequestObject struct {
	Params PostAuthLoginWithPasswordParams
	Body   *PostAuthLoginWithPasswordJSONRequestBody
}

type PostAuthLoginWithPasswordResponseObject interface {
	VisitPostAuthLoginWithPasswordResponse(ctx fiber.Ctx) error
}

type PostAuthLoginWithPassword200JSONResponse AuthenticatedUser

func (response PostAuthLoginWithPassword200JSONResponse) VisitPostAuthLoginWithPasswordResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type PostAuthLoginWithPassword400JSONResponse ValidationError

func (response PostAuthLoginWithPassword400JSONResponse) VisitPostAuthLoginWithPasswordResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type PostAuthLoginWithPassword404JSONResponse Err

func (response PostAuthLoginWithPassword404JSONResponse) VisitPostAuthLoginWithPasswordResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(404)

	return ctx.JSON(&response)
}

type PostAuthLogoutRequestObject struct {
	Body *PostAuthLogoutJSONRequestBody
}

type PostAuthLogoutResponseObject interface {
	VisitPostAuthLogoutResponse(ctx fiber.Ctx) error
}

type PostAuthLogout200Response struct {
}

func (response PostAuthLogout200Response) VisitPostAuthLogoutResponse(ctx fiber.Ctx) error {
	ctx.Status(200)
	return nil
}

type PostAuthLogout401JSONResponse Err

func (response PostAuthLogout401JSONResponse) VisitPostAuthLogoutResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type PostAuthRefreshRequestObject struct {
	Body *PostAuthRefreshJSONRequestBody
}

type PostAuthRefreshResponseObject interface {
	VisitPostAuthRefreshResponse(ctx fiber.Ctx) error
}

type PostAuthRefresh200JSONResponse Tokens

func (response PostAuthRefresh200JSONResponse) VisitPostAuthRefreshResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type PostAuthRefresh401JSONResponse Err

func (response PostAuthRefresh401JSONResponse) VisitPostAuthRefreshResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type PostAuthRegisterWithPasswordRequestObject struct {
	Params PostAuthRegisterWithPasswordParams
	Body   *PostAuthRegisterWithPasswordJSONRequestBody
}

type PostAuthRegisterWithPasswordResponseObject interface {
	VisitPostAuthRegisterWithPasswordResponse(ctx fiber.Ctx) error
}

type PostAuthRegisterWithPassword201JSONResponse AuthenticatedUser

func (response PostAuthRegisterWithPassword201JSONResponse) VisitPostAuthRegisterWithPasswordResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(201)

	return ctx.JSON(&response)
}

type PostAuthRegisterWithPassword400JSONResponse ValidationError

func (response PostAuthRegisterWithPassword400JSONResponse) VisitPostAuthRegisterWithPasswordResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type PostAuthRegisterWithPassword409JSONResponse Err

func (response PostAuthRegisterWithPassword409JSONResponse) VisitPostAuthRegisterWithPasswordResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(409)

	return ctx.JSON(&response)
}

type PostAuthVerifyEmailRequestObject struct {
	Params PostAuthVerifyEmailParams
}

type PostAuthVerifyEmailResponseObject interface {
	VisitPostAuthVerifyEmailResponse(ctx fiber.Ctx) error
}

type PostAuthVerifyEmail200Response struct {
}

func (response PostAuthVerifyEmail200Response) VisitPostAuthVerifyEmailResponse(ctx fiber.Ctx) error {
	ctx.Status(200)
	return nil
}

type PostAuthVerifyEmail400JSONResponse ValidationError

func (response PostAuthVerifyEmail400JSONResponse) VisitPostAuthVerifyEmailResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type GetHealthzRequestObject struct {
}

type GetHealthzResponseObject interface {
	VisitGetHealthzResponse(ctx fiber.Ctx) error
}

type GetHealthz200TextResponse string

func (response GetHealthz200TextResponse) VisitGetHealthzResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "text/plain")
	ctx.Status(200)

	_, err := ctx.WriteString(string(response))
	return err
}

type GetMetricsRequestObject struct {
}

type GetMetricsResponseObject interface {
	VisitGetMetricsResponse(ctx fiber.Ctx) error
}

type GetMetrics200JSONResponse map[string]interface{}

func (response GetMetrics200JSONResponse) VisitGetMetricsResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetPingRequestObject struct {
}

type GetPingResponseObject interface {
	VisitGetPingResponse(ctx fiber.Ctx) error
}

type GetPing200TextResponse string

func (response GetPing200TextResponse) VisitGetPingResponse(ctx fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "text/plain")
	ctx.Status(200)

	_, err := ctx.WriteString(string(response))
	return err
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Checks if access token is valid
	// (GET /auth/check)
	GetAuthCheck(ctx context.Context, request GetAuthCheckRequestObject) (GetAuthCheckResponseObject, error)
	// Login user by their password and email or nickname
	// (POST /auth/loginWithPassword)
	PostAuthLoginWithPassword(ctx context.Context, request PostAuthLoginWithPasswordRequestObject) (PostAuthLoginWithPasswordResponseObject, error)
	// Logouts user by invalidating refresh token
	// (POST /auth/logout)
	PostAuthLogout(ctx context.Context, request PostAuthLogoutRequestObject) (PostAuthLogoutResponseObject, error)
	// Refreshes expired access and unexpired refresh tokens
	// (POST /auth/refresh)
	PostAuthRefresh(ctx context.Context, request PostAuthRefreshRequestObject) (PostAuthRefreshResponseObject, error)
	// Registrates user by a password and email or nickname
	// (POST /auth/registerWithPassword)
	PostAuthRegisterWithPassword(ctx context.Context, request PostAuthRegisterWithPasswordRequestObject) (PostAuthRegisterWithPasswordResponseObject, error)
	// Verifies the user's email.
	// (POST /auth/verify/email)
	PostAuthVerifyEmail(ctx context.Context, request PostAuthVerifyEmailRequestObject) (PostAuthVerifyEmailResponseObject, error)
	// The server's health probes
	// (GET /healthz)
	GetHealthz(ctx context.Context, request GetHealthzRequestObject) (GetHealthzResponseObject, error)
	// OpenTelemetry metrics
	// (GET /metrics)
	GetMetrics(ctx context.Context, request GetMetricsRequestObject) (GetMetricsResponseObject, error)
	// Pings the server
	// (GET /ping)
	GetPing(ctx context.Context, request GetPingRequestObject) (GetPingResponseObject, error)
}
type StrictHandlerFunc func(ctx fiber.Ctx, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetAuthCheck operation middleware
func (sh *strictHandler) GetAuthCheck(ctx fiber.Ctx) error {
	var request GetAuthCheckRequestObject

	response, err := sh.ssi.GetAuthCheck(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitGetAuthCheckResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// PostAuthLoginWithPassword operation middleware
func (sh *strictHandler) PostAuthLoginWithPassword(ctx fiber.Ctx, params PostAuthLoginWithPasswordParams) error {
	var request PostAuthLoginWithPasswordRequestObject

	request.Params = params

	var body PostAuthLoginWithPasswordJSONRequestBody
	if err := ctx.Bind().JSON(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	// validation by https://github.com/Onnywrite
	if err := fmtvalidate.V.StructCtx(ctx.Context(), body); err != nil {
		return PostAuthLoginWithPassword400JSONResponse{
			Service: ValidationErrorServiceSsonny,
			Fields:  fmtvalidate.FormatFields(err),
		}.VisitPostAuthLoginWithPasswordResponse(ctx)
	}

	response, err := sh.ssi.PostAuthLoginWithPassword(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitPostAuthLoginWithPasswordResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// PostAuthLogout operation middleware
func (sh *strictHandler) PostAuthLogout(ctx fiber.Ctx) error {
	var request PostAuthLogoutRequestObject

	var body PostAuthLogoutJSONRequestBody
	if err := ctx.Bind().JSON(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	response, err := sh.ssi.PostAuthLogout(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitPostAuthLogoutResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// PostAuthRefresh operation middleware
func (sh *strictHandler) PostAuthRefresh(ctx fiber.Ctx) error {
	var request PostAuthRefreshRequestObject

	var body PostAuthRefreshJSONRequestBody
	if err := ctx.Bind().JSON(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	response, err := sh.ssi.PostAuthRefresh(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitPostAuthRefreshResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// PostAuthRegisterWithPassword operation middleware
func (sh *strictHandler) PostAuthRegisterWithPassword(ctx fiber.Ctx, params PostAuthRegisterWithPasswordParams) error {
	var request PostAuthRegisterWithPasswordRequestObject

	request.Params = params

	var body PostAuthRegisterWithPasswordJSONRequestBody
	if err := ctx.Bind().JSON(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	// validation by https://github.com/Onnywrite
	if err := fmtvalidate.V.StructCtx(ctx.Context(), body); err != nil {
		return PostAuthRegisterWithPassword400JSONResponse{
			Service: ValidationErrorServiceSsonny,
			Fields:  fmtvalidate.FormatFields(err),
		}.VisitPostAuthRegisterWithPasswordResponse(ctx)
	}

	response, err := sh.ssi.PostAuthRegisterWithPassword(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitPostAuthRegisterWithPasswordResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// PostAuthVerifyEmail operation middleware
func (sh *strictHandler) PostAuthVerifyEmail(ctx fiber.Ctx, params PostAuthVerifyEmailParams) error {
	var request PostAuthVerifyEmailRequestObject

	request.Params = params

	response, err := sh.ssi.PostAuthVerifyEmail(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitPostAuthVerifyEmailResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// GetHealthz operation middleware
func (sh *strictHandler) GetHealthz(ctx fiber.Ctx) error {
	var request GetHealthzRequestObject

	response, err := sh.ssi.GetHealthz(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitGetHealthzResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// GetMetrics operation middleware
func (sh *strictHandler) GetMetrics(ctx fiber.Ctx) error {
	var request GetMetricsRequestObject

	response, err := sh.ssi.GetMetrics(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitGetMetricsResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// GetPing operation middleware
func (sh *strictHandler) GetPing(ctx fiber.Ctx) error {
	var request GetPingRequestObject

	response, err := sh.ssi.GetPing(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := response.VisitGetPingResponse(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
