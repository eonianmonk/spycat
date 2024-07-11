package context

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	logKey = "logkeyctx"
)

func SetLogContext(c *fiber.Ctx, l log.Logger) {
	c.Locals(logKey, l)
}

func GetLogCtx(c *fiber.Ctx) log.Logger {
	return c.Locals(logKey).(log.Logger)
}
