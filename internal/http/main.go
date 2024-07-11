package http

import (
	"github.com/eonianmonk/spycat"
	"github.com/eonianmonk/spycat/internal/http/context"
	"github.com/gofiber/fiber/v2"
)

func Run(dbContext *context.DbsCtx, breedValidator spycat.Validator, endpoint string) {
	app := fiber.New()
	app.Use(ContentTypeMW())
	app.Use(func(c *fiber.Ctx) error {
		context.SetDbContext(c, dbContext)
		context.SetCatsBreedContext(c, breedValidator)
		return c.Next()
	})
	setupRoutes(app)
	err := app.Listen(endpoint)
	if err != nil {
		panic(err)
	}

}
