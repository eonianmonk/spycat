package handlers

import (
	"database/sql"

	"github.com/eonianmonk/spycat"
	"github.com/eonianmonk/spycat/internal/data"
	"github.com/eonianmonk/spycat/internal/http/context"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

func rollbackTxs(txs ...*sql.Tx) error {
	var err error
	for _, tx := range txs {
		err = tx.Rollback()
	}
	return err
}

func errWithRollback(err error, c *fiber.Ctx, txs ...*sql.Tx) error {
	log.Error(err)
	errTx := rollbackTxs(txs...)
	log.Error("failed to revert tx: ", errTx, txs)
	return c.SendStatus(fiber.StatusInternalServerError)
}

// creates mission, can create with targets
func CreateMissionWithTargets(c *fiber.Ctx) error {
	var missionReq data.Mission
	err := c.BodyParser(&missionReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	missionTx, errM := context.GetDbContext(c).MissionsDb.Db.Begin()
	targetsTx, errT := context.GetDbContext(c).TargetsDb.Db.Begin()
	if errM != nil || errT != nil {
		log.Error(errM, errT)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	mission, err := context.GetDbContext(c).MissionsDb.Create(&missionReq, missionTx)
	if err != nil {
		return errWithRollback(err, c, missionTx, targetsTx)
	}
	targets, err := context.GetDbContext(c).TargetsDb.CreateMany(mission.Targets, mission.Id, targetsTx)
	if err != nil {
		return errWithRollback(err, c, missionTx, targetsTx)
	}
	mission.Targets = targets
	// TODO: add logging on error
	targetsTx.Commit()
	missionTx.Commit()
	return c.JSON(mission)
}

func DeleteMission(c *fiber.Ctx) error {
	id, err := getIdParam(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).MissionsDb.Delete(id)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	log.Debug("Deleted mission %s", id)
	return c.SendStatus(fiber.StatusOK)
}

func UpdateMission(c *fiber.Ctx) error {
	var missionReq data.Mission
	err := c.BodyParser(&missionReq)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).MissionsDb.UpdateCompletion(missionReq.Id, spycat.Complete)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(missionReq)
}

func AssignCat(c *fiber.Ctx) error {
	catId, missionId, err := parseCatAssign(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = context.GetDbContext(c).MissionsDb.Assign(missionId, catId)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(data.Mission{Id: missionId, AssignedCatId: catId, Status: spycat.Incomplete})
}

func ListMissions(c *fiber.Ctx) error {
	offset, count, err := ParseListQuery(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	missions, err := context.GetDbContext(c).MissionsDb.List(offset, count)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = context.GetDbContext(c).TargetsDb.GetTargetsForMissions(missions)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(missions)
}

func GetMission(c *fiber.Ctx) error {
	id, err := getIdParam(c)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mission, err := context.GetDbContext(c).MissionsDb.Get(id)
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(mission)
}
