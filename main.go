package main

import (
	"strings"
	"wallet-topup-system/common/cache"
	"wallet-topup-system/common/logs"
	"wallet-topup-system/config"
	"wallet-topup-system/core/handlers"
	"wallet-topup-system/core/middlewares"
	"wallet-topup-system/core/repositories"
	"wallet-topup-system/core/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "wallet-topup-system/docs"
)

// @title           Wallet Topup System API
// @version         1.0.0
// @description     API สำหรับ Topup Wallet System
// @securityDefinitions.apikey ApiKeyAuth
// @in				header
// @name			X-API-KEY
// @contact.name    Julladith Klinloy
// @contact.email   julladith.kl@gmail.com

func init() {
	// initialize environment config
	config.NewAppInitEnvironment()

	// initialize timezone
	config.NewAppInitTime()
}

func main() {
	// initialize app
	db := config.NewAppInitDBPostgres()
	cache := cache.NewAppRedisCache()
	log := logs.NewAppZapLogs()

	// repositories
	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// services
	userSrv := services.NewUserService(log, userRepo)
	walletSrv := services.NewWalletService(log, cache, userRepo, transactionRepo)

	// handlers
	walletHandler := handlers.NewWalletHandler(walletSrv)
	userHandler := handlers.NewUserHandler(userSrv)

	// initialize user data
	userSrv.Initialize()

	// start app
	app := fiber.New(config.ServerConfig())

	// setup global middleware
	app.Use(recover.New())
	app.Use(cors.New(config.CorsConfig()))

	// routes
	app.Post("/wallet/verify", middlewares.APIKey(), walletHandler.TopUpVerified)
	app.Post("/wallet/confirm", middlewares.APIKey(), walletHandler.TopUpConfirmed)
	app.Get("/user/wallet/:user_id", middlewares.APIKey(), userHandler.GetUserWallet)

	// docs swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// listen server
	addr := ":" + config.Env.Port
	if env := strings.ToLower(config.Env.Env); env == "dev" || env == "development" {
		addr = "localhost" + addr
	}
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
