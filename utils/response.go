package utils

import "github.com/gofiber/fiber/v2"

// format using the Fiber framework
func ErrorFormat(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": msg,
	})
}

func BodyParserFail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusBadRequest,
		"message": "incorrect payload format",
	})
}

func ParamParserFail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusBadRequest,
		"message": "failed to parse params",
	})
}

func QueryParserFail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusBadRequest,
		"message": "failed to parse query params",
	})
}
