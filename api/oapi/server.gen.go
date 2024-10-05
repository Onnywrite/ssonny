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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Login user by their password and email or nickname
	// (POST /auth/loginWithPassword)
	PostAuthLoginWithPassword(c *fiber.Ctx, params PostAuthLoginWithPasswordParams) error
	// Login user by their password and email or nickname
	// (POST /auth/logout)
	PostAuthLogout(c *fiber.Ctx) error
	// Refreshes expired access and unexpired refresh tokens
	// (POST /auth/refresh)
	PostAuthRefresh(c *fiber.Ctx) error
	// Registrates user by a password and email or nickname
	// (POST /auth/registerWithPassword)
	PostAuthRegisterWithPassword(c *fiber.Ctx, params PostAuthRegisterWithPasswordParams) error
	// Verifies the user's email.
	// (POST /auth/verify/email)
	PostAuthVerifyEmail(c *fiber.Ctx, params PostAuthVerifyEmailParams) error
	// The server's health probes
	// (GET /healthz)
	GetHealthz(c *fiber.Ctx) error
	// OpenTelemetry metrics
	// (GET /metrics)
	GetMetrics(c *fiber.Ctx) error
	// Pings the server
	// (GET /ping)
	GetPing(c *fiber.Ctx) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// PostAuthLoginWithPassword operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLoginWithPassword(c *fiber.Ctx) error {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthLoginWithPasswordParams

	headers := c.GetReqHeaders()

	// ------------- Optional header parameter "User-Agent" -------------
	if values, found := headers[http.CanonicalHeaderKey("User-Agent")]; found {
		var UserAgent UserAgent

		err = runtime.BindStyledParameterWithOptions("simple", "User-Agent", values[0], &UserAgent, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter User-Agent: %w", err).Error())
		}

		params.UserAgent = &UserAgent

	}

	return siw.Handler.PostAuthLoginWithPassword(c, params)
}

// PostAuthLogout operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLogout(c *fiber.Ctx) error {

	return siw.Handler.PostAuthLogout(c)
}

// PostAuthRefresh operation middleware
func (siw *ServerInterfaceWrapper) PostAuthRefresh(c *fiber.Ctx) error {

	return siw.Handler.PostAuthRefresh(c)
}

// PostAuthRegisterWithPassword operation middleware
func (siw *ServerInterfaceWrapper) PostAuthRegisterWithPassword(c *fiber.Ctx) error {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthRegisterWithPasswordParams

	headers := c.GetReqHeaders()

	// ------------- Optional header parameter "User-Agent" -------------
	if values, found := headers[http.CanonicalHeaderKey("User-Agent")]; found {
		var UserAgent UserAgent

		err = runtime.BindStyledParameterWithOptions("simple", "User-Agent", values[0], &UserAgent, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter User-Agent: %w", err).Error())
		}

		params.UserAgent = &UserAgent

	}

	return siw.Handler.PostAuthRegisterWithPassword(c, params)
}

// PostAuthVerifyEmail operation middleware
func (siw *ServerInterfaceWrapper) PostAuthVerifyEmail(c *fiber.Ctx) error {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostAuthVerifyEmailParams

	var query url.Values
	query, err = url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "token" -------------

	if paramValue := c.Query("token"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument token is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "token", query, &params.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter token: %w", err).Error())
	}

	return siw.Handler.PostAuthVerifyEmail(c, params)
}

// GetHealthz operation middleware
func (siw *ServerInterfaceWrapper) GetHealthz(c *fiber.Ctx) error {

	return siw.Handler.GetHealthz(c)
}

// GetMetrics operation middleware
func (siw *ServerInterfaceWrapper) GetMetrics(c *fiber.Ctx) error {

	return siw.Handler.GetMetrics(c)
}

// GetPing operation middleware
func (siw *ServerInterfaceWrapper) GetPing(c *fiber.Ctx) error {

	return siw.Handler.GetPing(c)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(fiber.Handler(m))
	}

	router.Post(options.BaseURL+"/auth/loginWithPassword", wrapper.PostAuthLoginWithPassword)

	router.Post(options.BaseURL+"/auth/logout", wrapper.PostAuthLogout)

	router.Post(options.BaseURL+"/auth/refresh", wrapper.PostAuthRefresh)

	router.Post(options.BaseURL+"/auth/registerWithPassword", wrapper.PostAuthRegisterWithPassword)

	router.Post(options.BaseURL+"/auth/verify/email", wrapper.PostAuthVerifyEmail)

	router.Get(options.BaseURL+"/healthz", wrapper.GetHealthz)

	router.Get(options.BaseURL+"/metrics", wrapper.GetMetrics)

	router.Get(options.BaseURL+"/ping", wrapper.GetPing)

}

type PostAuthLoginWithPasswordRequestObject struct {
	Params PostAuthLoginWithPasswordParams
	Body   *PostAuthLoginWithPasswordJSONRequestBody
}

type PostAuthLoginWithPasswordResponseObject interface {
	VisitPostAuthLoginWithPasswordResponse(ctx *fiber.Ctx) error
}

type PostAuthLoginWithPassword200JSONResponse AuthenticatedUser

func (response PostAuthLoginWithPassword200JSONResponse) VisitPostAuthLoginWithPasswordResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type PostAuthLoginWithPassword400JSONResponse Err

func (response PostAuthLoginWithPassword400JSONResponse) VisitPostAuthLoginWithPasswordResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type PostAuthLoginWithPassword404JSONResponse Err

func (response PostAuthLoginWithPassword404JSONResponse) VisitPostAuthLoginWithPasswordResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(404)

	return ctx.JSON(&response)
}

type PostAuthLogoutRequestObject struct {
	Body *PostAuthLogoutJSONRequestBody
}

type PostAuthLogoutResponseObject interface {
	VisitPostAuthLogoutResponse(ctx *fiber.Ctx) error
}

type PostAuthLogout200Response struct {
}

func (response PostAuthLogout200Response) VisitPostAuthLogoutResponse(ctx *fiber.Ctx) error {
	ctx.Status(200)
	return nil
}

type PostAuthLogout401Response struct {
}

func (response PostAuthLogout401Response) VisitPostAuthLogoutResponse(ctx *fiber.Ctx) error {
	ctx.Status(401)
	return nil
}

type PostAuthRefreshRequestObject struct {
	Body *PostAuthRefreshJSONRequestBody
}

type PostAuthRefreshResponseObject interface {
	VisitPostAuthRefreshResponse(ctx *fiber.Ctx) error
}

type PostAuthRefresh200JSONResponse Tokens

func (response PostAuthRefresh200JSONResponse) VisitPostAuthRefreshResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type PostAuthRefresh401JSONResponse Err

func (response PostAuthRefresh401JSONResponse) VisitPostAuthRefreshResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type PostAuthRegisterWithPasswordRequestObject struct {
	Params PostAuthRegisterWithPasswordParams
	Body   *PostAuthRegisterWithPasswordJSONRequestBody
}

type PostAuthRegisterWithPasswordResponseObject interface {
	VisitPostAuthRegisterWithPasswordResponse(ctx *fiber.Ctx) error
}

type PostAuthRegisterWithPassword201JSONResponse AuthenticatedUser

func (response PostAuthRegisterWithPassword201JSONResponse) VisitPostAuthRegisterWithPasswordResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(201)

	return ctx.JSON(&response)
}

type PostAuthRegisterWithPassword400JSONResponse Err

func (response PostAuthRegisterWithPassword400JSONResponse) VisitPostAuthRegisterWithPasswordResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type PostAuthRegisterWithPassword409JSONResponse Err

func (response PostAuthRegisterWithPassword409JSONResponse) VisitPostAuthRegisterWithPasswordResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(409)

	return ctx.JSON(&response)
}

type PostAuthVerifyEmailRequestObject struct {
	Params PostAuthVerifyEmailParams
}

type PostAuthVerifyEmailResponseObject interface {
	VisitPostAuthVerifyEmailResponse(ctx *fiber.Ctx) error
}

type PostAuthVerifyEmail200Response struct {
}

func (response PostAuthVerifyEmail200Response) VisitPostAuthVerifyEmailResponse(ctx *fiber.Ctx) error {
	ctx.Status(200)
	return nil
}

type PostAuthVerifyEmail400JSONResponse Err

func (response PostAuthVerifyEmail400JSONResponse) VisitPostAuthVerifyEmailResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type GetHealthzRequestObject struct {
}

type GetHealthzResponseObject interface {
	VisitGetHealthzResponse(ctx *fiber.Ctx) error
}

type GetHealthz200TextResponse string

func (response GetHealthz200TextResponse) VisitGetHealthzResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "text/plain")
	ctx.Status(200)

	_, err := ctx.WriteString(string(response))
	return err
}

type GetMetricsRequestObject struct {
}

type GetMetricsResponseObject interface {
	VisitGetMetricsResponse(ctx *fiber.Ctx) error
}

type GetMetrics200JSONResponse map[string]interface{}

func (response GetMetrics200JSONResponse) VisitGetMetricsResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetPingRequestObject struct {
}

type GetPingResponseObject interface {
	VisitGetPingResponse(ctx *fiber.Ctx) error
}

type GetPing200TextResponse string

func (response GetPing200TextResponse) VisitGetPingResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "text/plain")
	ctx.Status(200)

	_, err := ctx.WriteString(string(response))
	return err
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Login user by their password and email or nickname
	// (POST /auth/loginWithPassword)
	PostAuthLoginWithPassword(ctx context.Context, request PostAuthLoginWithPasswordRequestObject) (PostAuthLoginWithPasswordResponseObject, error)
	// Login user by their password and email or nickname
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

type StrictHandlerFunc func(ctx *fiber.Ctx, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// PostAuthLoginWithPassword operation middleware
func (sh *strictHandler) PostAuthLoginWithPassword(ctx *fiber.Ctx, params PostAuthLoginWithPasswordParams) error {
	var request PostAuthLoginWithPasswordRequestObject

	request.Params = params

	var body PostAuthLoginWithPasswordJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.PostAuthLoginWithPassword(ctx.UserContext(), request.(PostAuthLoginWithPasswordRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAuthLoginWithPassword")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(PostAuthLoginWithPasswordResponseObject); ok {
		if err := validResponse.VisitPostAuthLoginWithPasswordResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// PostAuthLogout operation middleware
func (sh *strictHandler) PostAuthLogout(ctx *fiber.Ctx) error {
	var request PostAuthLogoutRequestObject

	var body PostAuthLogoutJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.PostAuthLogout(ctx.UserContext(), request.(PostAuthLogoutRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAuthLogout")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(PostAuthLogoutResponseObject); ok {
		if err := validResponse.VisitPostAuthLogoutResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// PostAuthRefresh operation middleware
func (sh *strictHandler) PostAuthRefresh(ctx *fiber.Ctx) error {
	var request PostAuthRefreshRequestObject

	var body PostAuthRefreshJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.PostAuthRefresh(ctx.UserContext(), request.(PostAuthRefreshRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAuthRefresh")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(PostAuthRefreshResponseObject); ok {
		if err := validResponse.VisitPostAuthRefreshResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// PostAuthRegisterWithPassword operation middleware
func (sh *strictHandler) PostAuthRegisterWithPassword(ctx *fiber.Ctx, params PostAuthRegisterWithPasswordParams) error {
	var request PostAuthRegisterWithPasswordRequestObject

	request.Params = params

	var body PostAuthRegisterWithPasswordJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.PostAuthRegisterWithPassword(ctx.UserContext(), request.(PostAuthRegisterWithPasswordRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAuthRegisterWithPassword")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(PostAuthRegisterWithPasswordResponseObject); ok {
		if err := validResponse.VisitPostAuthRegisterWithPasswordResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// PostAuthVerifyEmail operation middleware
func (sh *strictHandler) PostAuthVerifyEmail(ctx *fiber.Ctx, params PostAuthVerifyEmailParams) error {
	var request PostAuthVerifyEmailRequestObject

	request.Params = params

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.PostAuthVerifyEmail(ctx.UserContext(), request.(PostAuthVerifyEmailRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAuthVerifyEmail")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(PostAuthVerifyEmailResponseObject); ok {
		if err := validResponse.VisitPostAuthVerifyEmailResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetHealthz operation middleware
func (sh *strictHandler) GetHealthz(ctx *fiber.Ctx) error {
	var request GetHealthzRequestObject

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetHealthz(ctx.UserContext(), request.(GetHealthzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetHealthz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetHealthzResponseObject); ok {
		if err := validResponse.VisitGetHealthzResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetMetrics operation middleware
func (sh *strictHandler) GetMetrics(ctx *fiber.Ctx) error {
	var request GetMetricsRequestObject

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetMetrics(ctx.UserContext(), request.(GetMetricsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetMetrics")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetMetricsResponseObject); ok {
		if err := validResponse.VisitGetMetricsResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetPing operation middleware
func (sh *strictHandler) GetPing(ctx *fiber.Ctx) error {
	var request GetPingRequestObject

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetPing(ctx.UserContext(), request.(GetPingRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetPing")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetPingResponseObject); ok {
		if err := validResponse.VisitGetPingResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
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
