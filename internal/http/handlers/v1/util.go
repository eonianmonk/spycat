package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var (
	IdParameter        string = "id"
	CatIdParameter            = "cat"
	MissionIdParameter        = "mission"
	ListOffset                = "offset"
	ListLimit                 = "limit"
)

func getIdParam(c *fiber.Ctx) (string, error) {
	id := c.Params(IdParameter)
	if id == "" {
		return "", fmt.Errorf("no id parameter provided")
	}
	return id, nil
}

// gets offset and limit from c
// returns offset, limit,error
func ParseListQuery(c *fiber.Ctx) (int, int, error) {
	offsetStr := c.Query(ListOffset)
	limitStr := c.Query(ListLimit)
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert offset to integer: %s", err)
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert limit to integer: %s", err)
	}
	return offset, limit, nil
}

// returns catId,missionId
func parseCatAssign(c *fiber.Ctx) (string, string, error) {
	catId := c.Query(CatIdParameter)
	missionId := c.Query(MissionIdParameter)
	if catId == "" || missionId == "" {
		return "", "", fmt.Errorf("both catId(%s) and missionId(%s) required")
	}
	return catId, missionId, nil
}
