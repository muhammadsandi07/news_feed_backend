package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			if e, ok := err.(*AppError); ok {
				return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
			}

			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{"error": "internal server error"})
		},
	}
}
