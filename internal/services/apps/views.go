package apps

import (
	"github.com/go-playground/validator/v10"
)

type CreateAppData struct {
	Name        string   `validate:"required,max=64"`
	Description string   `validate:"omitempty,max=4096"`
	DomainsIds  []uint64 `validate:"omitempty,dive,fqdn"`
}

func (d CreateAppData) Validate(validate *validator.Validate) error {
	return validate.Struct(d)
}

type AppCreated struct {
	Id     uint64
	Secret string
}

type AppWithDomainsCreated struct {
	AppCreated
	DomainsIds []uint64
}
