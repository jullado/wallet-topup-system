package middlewares

import (
	"wallet-topup-system/config"
	"wallet-topup-system/core/models"
	"wallet-topup-system/utils"

	"github.com/gofiber/fiber/v2"
)

func APIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Env.APIKey != "" && c.Get("X-API-KEY") != config.Env.APIKey {
			return utils.ErrorFormat(c, fiber.StatusUnauthorized, models.ErrUnauthorized)
		}
		return c.Next()
	}
}
