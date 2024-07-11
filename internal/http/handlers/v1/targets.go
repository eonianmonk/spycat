package handlers

import (
	"github.com/eonianmonk/spycat/internal/data"
	"github.com/eonianmonk/spycat/internal/http/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateTarget(c *fiber.Ctx) error {
	targetReq := data.Target{}
	err := c.BodyParser(&targetReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	targets, err := context.GetDbContext(c).TargetsDb.CreateMany([]*data.Target{&targetReq}, targetReq.MissionId, nil)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(targets[0])
}

func DeleteTarget(c *fiber.Ctx) error {
	id, err := getIdParam(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).TargetsDb.Delete(id)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	log.Debug("Deleted target %s", id)
	return c.SendStatus(fiber.StatusOK)
}

func UpdateTarget(c *fiber.Ctx) error {
	targetReq := data.Target{}
	log.Error()
	err := c.BodyParser(&targetReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).TargetsDb.Update(targetReq.Id, targetReq.Notes, targetReq.Status)
	if err != nil {
		log.Error(targetReq.Id, err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(targetReq)
}
