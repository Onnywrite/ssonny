// ServerInterface represents all server handlers.
type ServerInterface interface {
{{range .}}{{.SummaryAsComment }}
// ({{.Method}} {{.Path}})
{{.OperationId}}(c fiber.Ctx{{genParamArgs .PathParams}}{{if .RequiresParamObject}}, params {{.OperationId}}Params{{end}}) error
{{end}}
}