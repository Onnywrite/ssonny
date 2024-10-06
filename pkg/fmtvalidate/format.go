package fmtvalidate

import (
	"errors"

	locales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	V     = validator.New(validator.WithRequiredStructEnabled())
	trans ut.Translator
)

func init() {
	locale := locales.New()
	uni := ut.New(locale, locale)

	trans, _ = uni.GetTranslator("en")
	translations.RegisterDefaultTranslations(V, trans)
}

func FormatFields(err error) map[string]any {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		fields := make(map[string]any, len(ve))

		for _, e := range ve {
			fields[e.Field()] = e.Translate(trans)
		}

		return fields
	}

	return map[string]any{
		"<unexpected_error>": err.Error(),
	}
}
