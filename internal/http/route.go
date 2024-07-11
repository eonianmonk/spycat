package http

import (
	"fmt"

	"github.com/eonianmonk/spycat/internal/http/handlers/v1"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/v1")

	cats := api.Group("/cats")
	cats.Post("", handlers.CreateCat)
	cats.Delete(fmt.Sprintf("/:%s", handlers.IdParameter), handlers.DeleteCat)
	cats.Patch("", handlers.UpdateCat)
	cats.Get("", handlers.ListCats)
	cats.Get(fmt.Sprintf("/:%s", handlers.IdParameter), handlers.GetCat)
	cats.Post(fmt.Sprintf("/:%s/assign/:%s", handlers.CatIdParameter, handlers.MissionIdParameter), handlers.AssignCat)

	missions := api.Group("/missions")
	missions.Post("", handlers.CreateMissionWithTargets)
	missions.Delete(fmt.Sprintf("/:%s", handlers.IdParameter), handlers.DeleteMission)
	missions.Patch("", handlers.UpdateMission)
	missions.Get("", handlers.ListMissions)
	missions.Get(fmt.Sprintf("/:%s", handlers.IdParameter), handlers.GetMission)

	targets := api.Group("/targets")
	targets.Post("", handlers.CreateTarget)
	targets.Delete(fmt.Sprintf("/:%s", handlers.IdParameter), handlers.DeleteTarget)
	targets.Patch("", handlers.UpdateTarget)
}
