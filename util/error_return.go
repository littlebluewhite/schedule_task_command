package util

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func Err(c *fiber.Ctx, err error, inner int) error {
	return c.Status(484).JSON(fiber.Map{"message": err.Error(), "inner_code": inner})
}

type JsonErr interface {
	Error() string
	MarshalJSON() ([]byte, error)
}

type MyErr string

func (e MyErr) Error() string {
	return string(e)
}

func (e MyErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}
