package handlers_api_auth

import (
	"context"
	"strings"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/fiberutil"
	"github.com/Onnywrite/ssonny/internal/services/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/mileusna/useragent"
)

type Registrator interface {
	RegisterWithPassword(ctx context.Context, data auth.RegisterWithPasswordData) (*auth.AuthenticatedUser, error)
}

func RegisterWithPassword(service Registrator) func(c fiber.Ctx) error {
	type registerData struct {
		Nickname *string
		Email    string
		Gender   *string
		Birthday *string
		Password string
	}

	return func(c fiber.Ctx) error {
		var data registerData
		if err := c.Bind().JSON(&data); err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		ua := useragent.Parse(c.Get("User-Agent"))

		platform := strings.Join([]string{ua.OS, ua.OSVersion}, " ")
		agent := strings.Join([]string{ua.Name, ua.Version}, " ")

		var birthday *time.Time
		if data.Birthday != nil {
			b, err := time.Parse(time.DateOnly, *data.Birthday)
			if err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}
			birthday = &b
		}
		authUser, err := service.RegisterWithPassword(c.Context(), auth.RegisterWithPasswordData{
			Nickname: data.Nickname,
			Email:    data.Email,
			Gender:   data.Gender,
			Birthday: birthday,
			Password: data.Password,
			UserInfo: auth.UserInfo{
				Platform: platform,
				Agent:    agent,
			},
		})
		if err != nil {
			return fiberutil.Error(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(authUser)
	}
}
