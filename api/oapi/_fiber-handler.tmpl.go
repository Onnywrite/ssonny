// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
    BaseURL string
    Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
  RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
{{if .}}wrapper := ServerInterfaceWrapper{
Handler: si,
}

for _, m := range options.Middlewares {
    router.Use(fiber.Handler(m))
}
{{end}}
{{range .}}
router.{{.Method | lower | title }}(options.BaseURL+"{{.Path | swaggerUriToFiberUri}}", wrapper.{{.OperationId}})
{{end}}
}