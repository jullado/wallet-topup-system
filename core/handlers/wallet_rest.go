package handlers

import (
	"wallet-topup-system/core/models"
	"wallet-topup-system/core/services"
	"wallet-topup-system/utils"

	"github.com/gofiber/fiber/v2"
)

type walletHandler struct {
	walletSrv services.WalletService
}

func NewWalletHandler(walletSrv services.WalletService) walletHandler {
	return walletHandler{
		walletSrv,
	}
}

// @Summary      Verify a Top-up Transaction
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Security 	 ApiKeyAuth
// @Param        request body models.HandTopUpVerifiedReqModel true "payload สำหรับสร้าง Top-up Transaction โดยมีเวลาหมดอายุ 1 นาที"
// @Success      200 {object} models.SrvTopUpVerifiedResModel
// @Failure      400 {object} utils.ErrHandler
// @Router       /wallet/verify [post]
func (h walletHandler) TopUpVerified(c *fiber.Ctx) error {
	// parse body
	body := models.HandTopUpVerifiedReqModel{}
	if err := c.BodyParser(&body); err != nil {
		return utils.BodyParserFail(c)
	}

	// call service
	res, err := h.walletSrv.TopUpVerified(body.UserID, body.Amount, body.PaymentMethod)
	if err != nil {
		appErr, ok := err.(utils.ErrHandler)
		if ok {
			return utils.ErrorFormat(c, appErr.Code, appErr.Message)
		}

		return utils.ErrorFormat(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// @Summary      Confirm a Top-up Transaction
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Security 	 ApiKeyAuth
// @Param        request body models.HandTopUpConfirmedReqModel true "payload สำหรับ confirm Top-up Transaction"
// @Success      200 {object} models.SrvTopUpConfirmedResModel
// @Failure      400 {object} utils.ErrHandler
// @Router       /wallet/confirm [post]
func (h walletHandler) TopUpConfirmed(c *fiber.Ctx) error {
	// parse body
	body := models.HandTopUpConfirmedReqModel{}
	if err := c.BodyParser(&body); err != nil {
		return utils.BodyParserFail(c)
	}

	// call service
	res, err := h.walletSrv.TopUpConfirmed(body.TransactionID)
	if err != nil {
		appErr, ok := err.(utils.ErrHandler)
		if ok {
			return utils.ErrorFormat(c, appErr.Code, appErr.Message)
		}

		return utils.ErrorFormat(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
