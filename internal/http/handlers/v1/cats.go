package handlers

import (
	"github.com/eonianmonk/spycat/internal/http/context"

	"github.com/eonianmonk/spycat/internal/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateCat(c *fiber.Ctx) error {
	catReq := data.Cat{}
	err := c.BodyParser(&catReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = context.GetCatsBreedContext(c).Validate(catReq.Breed)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	cat, err := context.GetDbContext(c).CatsDb.Create(&catReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(cat)
}

func DeleteCat(c *fiber.Ctx) error {
	id, err := getIdParam(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).CatsDb.Delete(id)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	log.Debug("Deleted cat %s", id)
	return c.SendStatus(fiber.StatusOK)
}

func UpdateCat(c *fiber.Ctx) error {
	catReq := data.Cat{}
	err := c.BodyParser(&catReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).CatsDb.UpdateSalary(catReq.Id, catReq.Salary)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(catReq)
}

func ListCats(c *fiber.Ctx) error {
	offset, limit, err := ParseListQuery(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	cats, err := context.GetDbContext(c).CatsDb.List(offset, limit)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(cats)
}

func GetCat(c *fiber.Ctx) error {
	id, err := getIdParam(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	cat, err := context.GetDbContext(c).CatsDb.GetCat(id)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(cat)
}
