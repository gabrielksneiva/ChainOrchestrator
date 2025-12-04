package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler middleware global de tratamento de erros
func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		logger.Error("request error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.Int("status", code),
		)

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
			"code":  code,
		})
	}
}
