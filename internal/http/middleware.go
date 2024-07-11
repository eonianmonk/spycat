package http

import "github.com/gofiber/fiber/v2"

const ()

// adds application/json content type to reqs
func ContentTypeMW() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get(fiber.HeaderContentType) == "" {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		}
		return c.Next()
	}
}
