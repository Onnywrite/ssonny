// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
    Middlewares []fiber.Handler
    EndpointMiddlewares map[string][]fiber.Handler
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
  RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

{{if .}}
// created by github.com/Onnywrite
// Constants for all endpoints
const (
  {{ range . -}}
  // {{.Method}} {{.Path}}: {{.Summary }}
  EP_{{.OperationId}} = "{{.Path | swaggerUriToFiberUri}}"
  {{end}}
)
{{end}}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
{{if .}}wrapper := ServerInterfaceWrapper{
Handler: si,
}

for _, m := range options.Middlewares {
    router.Use(m)
}

mws := func (ep string) []fiber.Handler {
  return options.EndpointMiddlewares[ep]
}
{{end}}
{{range . -}}
router.{{.Method | lower | title }}(EP_{{.OperationId}}, wrapper.{{.OperationId}}, mws(EP_{{.OperationId}})...)
{{end}}
}