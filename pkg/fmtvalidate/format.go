package fmtvalidate

import (
	"errors"

	locales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	//nolint: varnamelen
	V     = validator.New(validator.WithRequiredStructEnabled())
	trans ut.Translator
)

// nolint: gochecknoinits
func init() {
	locale := locales.New()
	uni := ut.New(locale, locale)

	var found bool

	trans, found = uni.GetTranslator("en")
	if !found {
		panic("fmtvalidate: translator not found")
	}

	err := translations.RegisterDefaultTranslations(V, trans)
	if err != nil {
		panic(err)
	}
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
