type StrictHandlerFunc func(ctx fiber.Ctx, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
    return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
    ssi StrictServerInterface
    middlewares []StrictMiddlewareFunc
}

{{ define "validate" -}}
{{$ucopid := ucFirst .OperationId}}
{{ range .Responses }}
{{ if eq .StatusCode "400"}}
// validation by https://github.com/Onnywrite
if err := fmtvalidate.V.StructCtx(ctx.Context(), body); err != nil {
    {{$typeName := printf "%s%s%s%s" $ucopid .StatusCode (index .Contents 0).NameTagOrContentType "Response" -}}
    return {{$typeName}}{
        Service:      ValidationErrorServiceSsonny,
        Fields:       fmtvalidate.FormatFields(err),
    }.Visit{{$ucopid}}Response(ctx)
}
{{ end }}
{{ end }}
{{ end -}}
{{range .}}
{{$opid := .OperationId}}
{{$operation := .}}
    // {{$opid}} operation middleware
    func (sh *strictHandler) {{.OperationId}}(ctx fiber.Ctx{{genParamArgs .PathParams}}{{if .RequiresParamObject}}, params {{.OperationId}}Params{{end}}) error {
        var request {{$opid | ucFirst}}RequestObject

        {{range .PathParams -}}
            {{$varName := .GoVariableName -}}
            request.{{.GoName}} = {{.GoVariableName}}
        {{end -}}

        {{if .RequiresParamObject -}}
            request.Params = params
        {{end -}}

        {{ if .HasMaskedRequestContentTypes -}}
            request.ContentType = string(ctx.Request().Header.ContentType())
        {{end -}}

        {{$multipleBodies := gt (len .Bodies) 1 -}}
        {{range .Bodies -}}
            {{if $multipleBodies}}if strings.HasPrefix(string(ctx.Request().Header.ContentType()), "{{.ContentType}}") { {{end}}
                {{if .IsJSON }}
                    var body {{$opid}}{{.NameTag}}RequestBody
                    if err := ctx.Bind().JSON(&body); err != nil {
                        return fiber.NewError(fiber.StatusBadRequest, err.Error())
                    }
                    request.{{if $multipleBodies}}{{.NameTag}}{{end}}Body = &body
					
                    {{ template "validate" $operation }}
				{{else if eq .NameTag "Formdata" }}
                    var body {{$opid}}{{.NameTag}}RequestBody
                    if err := ctx.Body().Form(&body); err != nil {
                        return fiber.NewError(fiber.StatusBadRequest, err.Error())
                    }
                    request.{{if $multipleBodies}}{{.NameTag}}{{end}}Body = &body

					{{ template "validate" $operation }}
                {{else if eq .NameTag "Multipart" -}}
                    {{if eq .ContentType "multipart/form-data" -}}
                    request.{{if $multipleBodies}}{{.NameTag}}{{end}}Body = multipart.NewReader(bytes.NewReader(ctx.Request().Body()), string(ctx.Request().Header.MultipartFormBoundary()))
                    {{else -}}
                    if _, params, err := mime.ParseMediaType(string(ctx.Request().Header.ContentType())); err != nil {
                        return fiber.NewError(fiber.StatusBadRequest, err.Error())
                    } else if boundary := params["boundary"]; boundary == "" {
                        return fiber.NewError(fiber.StatusBadRequest, http.ErrMissingBoundary.Error())
                    } else {
                        request.{{if $multipleBodies}}{{.NameTag}}{{end}}Body = multipart.NewReader(bytes.NewReader(ctx.Request().Body()), boundary)
                    }
                    {{end -}}
                {{else if eq .NameTag "Text" -}}
                    data := ctx.Request().Body()
                    body := {{$opid}}{{.NameTag}}RequestBody(data)
                    request.{{if $multipleBodies}}{{.NameTag}}{{end}}Body = &body

					{{ template "validate" $operation }}
                {{else -}}
                    request.{{if $multipleBodies}}{{.NameTag}}{{end}}Body = bytes.NewReader(ctx.Request().Body())
                {{end}}{{/* if eq .NameTag "JSON" */ -}}
            {{if $multipleBodies}}}{{end}}
        {{end}}{{/* range .Bodies */}}

        response, err := sh.ssi.{{.OperationId}}(ctx.UserContext(), request)
        if err != nil {
            return fiber.NewError(fiber.StatusInternalServerError, err.Error())
        }
        
        if err := response.Visit{{$opid}}Response(ctx); err != nil {
            return fiber.NewError(fiber.StatusInternalServerError, err.Error())
        }

        return nil
    }
{{end}}