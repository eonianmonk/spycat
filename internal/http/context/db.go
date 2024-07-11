package context

import (
	"github.com/eonianmonk/spycat/internal/data"

	"github.com/gofiber/fiber/v2"
)

const (
	dbMiddlewareKey = "dbkeyctx"
)

type DbsCtx struct {
	CatsDb     *data.CatsDb
	MissionsDb *data.MissionsDb
	TargetsDb  *data.TargetDb
}

func SetDbContext(c *fiber.Ctx, mw *DbsCtx) {
	c.Locals(dbMiddlewareKey, mw)
}

func GetDbContext(c *fiber.Ctx) *DbsCtx {
	return c.Locals(dbMiddlewareKey).(*DbsCtx)
}
