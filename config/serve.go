package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ServerConfig() fiber.Config {
	return fiber.Config{
		ServerHeader:          "Wallet Topup System",
		AppName:               "Wallet Topup System",
		DisableStartupMessage: false,
		ReadTimeout:           Env.ReadTimeout,
		WriteTimeout:          Env.WriteTimeout,
	}
}

func CorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: Env.Cors,
		AllowHeaders: "Origin, Content-Type, Accept, X-API-KEY",
	}
}
