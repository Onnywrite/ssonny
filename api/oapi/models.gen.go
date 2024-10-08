// Package httpapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package httpapi

import (
	"time"

	googleuuid "github.com/google/uuid"
)

// Defines values for ErrService.
const (
	ErrServiceSsonny ErrService = "ssonny"
)

// Defines values for ValidationErrorService.
const (
	ValidationErrorServiceSsonny ValidationErrorService = "ssonny"
)

// AuthenticatedUser defines model for AuthenticatedUser.
type AuthenticatedUser struct {
	Access  string  `json:"Access"`
	Profile Profile `json:"Profile"`
	Refresh string  `json:"Refresh"`
}

// Err defines model for Err.
type Err struct {
	Message string     `json:"Message"`
	Service ErrService `json:"Service"`
}

// ErrService defines model for Err.Service.
type ErrService string

// Profile defines model for Profile.
type Profile struct {
	Birthday  *string         `json:"Birthday,omitempty"`
	CreatedAt time.Time       `json:"CreatedAt"`
	Email     string          `json:"Email"`
	Gender    *string         `json:"Gender,omitempty"`
	Id        googleuuid.UUID `json:"Id"`
	Nickname  *string         `json:"Nickname,omitempty"`
}

// Tokens defines model for Tokens.
type Tokens struct {
	Access  string `json:"Access"`
	Refresh string `json:"Refresh"`
}

// ValidationError defines model for ValidationError.
type ValidationError struct {
	Fields  map[string]interface{} `json:"Fields"`
	Service ValidationErrorService `json:"Service"`
}

// ValidationErrorService defines model for ValidationError.Service.
type ValidationErrorService string

// EmailToken defines model for EmailToken.
type EmailToken = string

// UserAgent defines model for UserAgent.
type UserAgent = string

// LoginWithPassword defines model for LoginWithPassword.
type LoginWithPassword struct {
	Email    *string `json:"Email,omitempty" validate:"omitempty,email,max=345"`
	Nickname *string `json:"Nickname,omitempty" validate:"omitempty,min=3,max=32"`
	Password string  `json:"Password" validate:"min=8,max=72"`
}

// Refresh defines model for Refresh.
type Refresh struct {
	RefreshToken string `json:"RefreshToken"`
}

// RegisterWithPassword defines model for RegisterWithPassword.
type RegisterWithPassword struct {
	Birthday *string `json:"Birthday,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Email    string  `json:"Email" validate:"email,max=345"`
	Gender   *string `json:"Gender,omitempty" validate:"omitempty,max=16"`
	Nickname *string `json:"Nickname,omitempty" validate:"omitempty,min=3,max=32"`
	Password string  `json:"Password" validate:"min=8,max=72"`
}

// PostAuthLoginWithPasswordJSONBody defines parameters for PostAuthLoginWithPassword.
type PostAuthLoginWithPasswordJSONBody struct {
	Email    *string `json:"Email,omitempty" validate:"omitempty,email,max=345"`
	Nickname *string `json:"Nickname,omitempty" validate:"omitempty,min=3,max=32"`
	Password string  `json:"Password" validate:"min=8,max=72"`
}

// PostAuthLoginWithPasswordParams defines parameters for PostAuthLoginWithPassword.
type PostAuthLoginWithPasswordParams struct {
	// UserAgent Default HTTP header in almost all browsers, don't care about this
	UserAgent UserAgent `json:"User-Agent"`
}

// PostAuthLogoutJSONBody defines parameters for PostAuthLogout.
type PostAuthLogoutJSONBody struct {
	RefreshToken string `json:"RefreshToken"`
}

// PostAuthRefreshJSONBody defines parameters for PostAuthRefresh.
type PostAuthRefreshJSONBody struct {
	RefreshToken string `json:"RefreshToken"`
}

// PostAuthRegisterWithPasswordJSONBody defines parameters for PostAuthRegisterWithPassword.
type PostAuthRegisterWithPasswordJSONBody struct {
	Birthday *string `json:"Birthday,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Email    string  `json:"Email" validate:"email,max=345"`
	Gender   *string `json:"Gender,omitempty" validate:"omitempty,max=16"`
	Nickname *string `json:"Nickname,omitempty" validate:"omitempty,min=3,max=32"`
	Password string  `json:"Password" validate:"min=8,max=72"`
}

// PostAuthRegisterWithPasswordParams defines parameters for PostAuthRegisterWithPassword.
type PostAuthRegisterWithPasswordParams struct {
	// UserAgent Default HTTP header in almost all browsers, don't care about this
	UserAgent UserAgent `json:"User-Agent"`
}

// PostAuthVerifyEmailParams defines parameters for PostAuthVerifyEmail.
type PostAuthVerifyEmailParams struct {
	// Token Special token for email verification
	Token EmailToken `form:"token" json:"token"`
}

// PostAuthLoginWithPasswordJSONRequestBody defines body for PostAuthLoginWithPassword for application/json ContentType.
type PostAuthLoginWithPasswordJSONRequestBody PostAuthLoginWithPasswordJSONBody

// PostAuthLogoutJSONRequestBody defines body for PostAuthLogout for application/json ContentType.
type PostAuthLogoutJSONRequestBody PostAuthLogoutJSONBody

// PostAuthRefreshJSONRequestBody defines body for PostAuthRefresh for application/json ContentType.
type PostAuthRefreshJSONRequestBody PostAuthRefreshJSONBody

// PostAuthRegisterWithPasswordJSONRequestBody defines body for PostAuthRegisterWithPassword for application/json ContentType.
type PostAuthRegisterWithPasswordJSONRequestBody PostAuthRegisterWithPasswordJSONBody
