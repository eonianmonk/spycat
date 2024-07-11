package context

import (
	"github.com/eonianmonk/spycat"
	"github.com/gofiber/fiber/v2"
)

const (
	catsBreedMiddlewareKey = "catsBreedKey"
)

func SetCatsBreedContext(c *fiber.Ctx, v spycat.Validator) {
	c.Locals(catsBreedMiddlewareKey, v)
}

func GetCatsBreedContext(c *fiber.Ctx) spycat.Validator {
	return c.Locals(catsBreedMiddlewareKey).(spycat.Validator)
}
