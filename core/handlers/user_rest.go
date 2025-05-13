package handlers

import (
	"wallet-topup-system/core/models"
	"wallet-topup-system/core/services"
	"wallet-topup-system/utils"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) userHandler {
	return userHandler{userSrv}
}

// @Summary      Get user wallet
// @Description  สําหรับดึงข้อมูล wallet ของ user โดยกำหนด user_id เริ่มต้นตอนสร้างระบบ คือ 1
// @Tags         User
// @Accept       json
// @Produce      json
// @Security 	 ApiKeyAuth
// @Param        user_id path uint true "user_id สําหรับดึงข้อมูล wallet ของ user"
// @Success      200 {object} models.SrvUserWalletModel
// @Failure      400 {object} utils.ErrHandler
// @Router       /user/wallet/{user_id} [get]
func (h *userHandler) GetUserWallet(c *fiber.Ctx) error {
	// parse params
	params := models.HandUserWalletReqModel{}
	if err := c.ParamsParser(&params); err != nil {
		return utils.ParamParserFail(c)
	}

	// call service
	res, err := h.userSrv.GetUserWallet(params.UserID)
	if err != nil {
		appErr, ok := err.(utils.ErrHandler)
		if ok {
			return utils.ErrorFormat(c, appErr.Code, appErr.Message)
		}

		return utils.ErrorFormat(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
