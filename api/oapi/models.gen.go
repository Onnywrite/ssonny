// Package httpapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package httpapi

import (
	"encoding/json"
	"time"

	googleuuid "github.com/google/uuid"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for ErrService.
const (
	Ssonny ErrService = "ssonny"
)

// AuthenticatedUser defines model for AuthenticatedUser.
type AuthenticatedUser struct {
	Access  string  `json:"Access"`
	Profile Profile `json:"Profile"`
	Refresh string  `json:"Refresh"`
}

// Err defines model for Err.
type Err struct {
	ErrorMessage string     `json:"ErrorMessage"`
	Service      ErrService `json:"Service"`
}

// ErrService defines model for Err.Service.
type ErrService string

// Profile defines model for Profile.
type Profile struct {
	Birthday  string              `json:"Birthday"`
	CreatedAt time.Time           `json:"CreatedAt"`
	Email     openapi_types.Email `json:"Email"`
	Gender    string              `json:"Gender"`
	Id        googleuuid.UUID     `json:"Id"`
	Nickname  string              `json:"Nickname"`
}

// RequestLoginWithPasswordAndEmail defines model for RequestLoginWithPasswordAndEmail.
type RequestLoginWithPasswordAndEmail struct {
	Email    openapi_types.Email `json:"Email" validate:"email,max=345"`
	Password string              `json:"Password" validate:"min=8,max=72"`
}

// RequestLoginWithPasswordAndNickname defines model for RequestLoginWithPasswordAndNickname.
type RequestLoginWithPasswordAndNickname struct {
	Nickname string `json:"Nickname" validate:"omitempty,min=3,max=32"`
	Password string `json:"Password" validate:"min=8,max=72"`
}

// Tokens defines model for Tokens.
type Tokens struct {
	Access  string `json:"Access"`
	Refresh string `json:"Refresh"`
}

// EmailToken defines model for EmailToken.
type EmailToken = string

// UserAgent defines model for UserAgent.
type UserAgent = string

// LoginWithPassword defines model for LoginWithPassword.
type LoginWithPassword struct {
	union json.RawMessage
}

// Refresh defines model for Refresh.
type Refresh struct {
	RefreshToken string `json:"RefreshToken"`
}

// RegisterWithPassword defines model for RegisterWithPassword.
type RegisterWithPassword struct {
	Birthday *string             `json:"Birthday,omitempty" validate:"omitempty,date"`
	Email    openapi_types.Email `json:"Email" validate:"email,max=345"`
	Gender   *string             `json:"Gender,omitempty" validate:"omitempty,max=16"`
	Nickname *string             `json:"Nickname,omitempty" validate:"omitempty,min=3,max=32"`
	Password string              `json:"Password" validate:"min=8,max=72"`
}

// PostAuthLoginWithPasswordJSONBody defines parameters for PostAuthLoginWithPassword.
type PostAuthLoginWithPasswordJSONBody struct {
	union json.RawMessage
}

// PostAuthLoginWithPasswordParams defines parameters for PostAuthLoginWithPassword.
type PostAuthLoginWithPasswordParams struct {
	// UserAgent Default HTTP header in almost all browsers, don't care about this
	UserAgent *UserAgent `json:"User-Agent,omitempty"`
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
	Birthday *string             `json:"Birthday,omitempty" validate:"omitempty,date"`
	Email    openapi_types.Email `json:"Email" validate:"email,max=345"`
	Gender   *string             `json:"Gender,omitempty" validate:"omitempty,max=16"`
	Nickname *string             `json:"Nickname,omitempty" validate:"omitempty,min=3,max=32"`
	Password string              `json:"Password" validate:"min=8,max=72"`
}

// PostAuthRegisterWithPasswordParams defines parameters for PostAuthRegisterWithPassword.
type PostAuthRegisterWithPasswordParams struct {
	// UserAgent Default HTTP header in almost all browsers, don't care about this
	UserAgent *UserAgent `json:"User-Agent,omitempty"`
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

// AsRequestLoginWithPasswordAndEmail returns the union data inside the LoginWithPassword as a RequestLoginWithPasswordAndEmail
func (t LoginWithPassword) AsRequestLoginWithPasswordAndEmail() (RequestLoginWithPasswordAndEmail, error) {
	var body RequestLoginWithPasswordAndEmail
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromRequestLoginWithPasswordAndEmail overwrites any union data inside the LoginWithPassword as the provided RequestLoginWithPasswordAndEmail
func (t *LoginWithPassword) FromRequestLoginWithPasswordAndEmail(v RequestLoginWithPasswordAndEmail) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeRequestLoginWithPasswordAndEmail performs a merge with any union data inside the LoginWithPassword, using the provided RequestLoginWithPasswordAndEmail
func (t *LoginWithPassword) MergeRequestLoginWithPasswordAndEmail(v RequestLoginWithPasswordAndEmail) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsRequestLoginWithPasswordAndNickname returns the union data inside the LoginWithPassword as a RequestLoginWithPasswordAndNickname
func (t LoginWithPassword) AsRequestLoginWithPasswordAndNickname() (RequestLoginWithPasswordAndNickname, error) {
	var body RequestLoginWithPasswordAndNickname
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromRequestLoginWithPasswordAndNickname overwrites any union data inside the LoginWithPassword as the provided RequestLoginWithPasswordAndNickname
func (t *LoginWithPassword) FromRequestLoginWithPasswordAndNickname(v RequestLoginWithPasswordAndNickname) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeRequestLoginWithPasswordAndNickname performs a merge with any union data inside the LoginWithPassword, using the provided RequestLoginWithPasswordAndNickname
func (t *LoginWithPassword) MergeRequestLoginWithPasswordAndNickname(v RequestLoginWithPasswordAndNickname) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t LoginWithPassword) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *LoginWithPassword) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
