package util

import "github.com/gofiber/fiber/v2"

func Err(c *fiber.Ctx, err error, inner int) error {
	return c.Status(484).JSON(fiber.Map{"message": err.Error(), "inner_code": inner})
}
