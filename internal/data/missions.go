package data

import (
	"database/sql"
	"fmt"

	"github.com/eonianmonk/spycat"
	"github.com/google/uuid"
)

type Mission struct {
	Id            string                 `json:"id,omitempty"`
	AssignedCatId string                 `json:"assigned_cat_id,omitempty"`
	Status        spycat.ComletionStatus `json:"completion_status,omitempty"`
	Targets       []*Target              `json:"targets,omitempty"`
}

type MissionsDb struct {
	Db *sql.DB
}

// function creates a mission in passed tx
func (mdb *MissionsDb) Create(mission *Mission, tx *sql.Tx) (*Mission, error) {
	assignedCatId := uuid.NullUUID{Valid: false}
	var err error
	if mission.AssignedCatId != "" {
		assignedCatId.UUID, err = uuid.Parse(mission.AssignedCatId)
		if err != nil {
			return nil, fmt.Errorf("invalid assigned cat uuid: %s", err)
		}
		assignedCatId.Valid = true
	}
	err = tx.QueryRow(`
		insert into missions (
			assigned_cat_id, completion_status
		) values ($1,$2)
		RETURNING id
	`, assignedCatId, mission.Status).Scan(&mission.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to create mission: %s", err)
	}
	return mission, nil
}

func (mdb *MissionsDb) Delete(uuid string) error {
	_, err := mdb.Db.Query(`DELETE FROM missions where id = $1`, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete mission %s: %s", uuid, err)
	}
	return nil
}

// updates mission in db - status, assigned cat
func (mdb *MissionsDb) UpdateCompletion(uuid string, status spycat.ComletionStatus) error {
	_, err := mdb.Db.Exec(`
		UPDATE missions SET
			completion_status = $2
		WHERE id = $1
	`, uuid, status)
	if err != nil {
		return fmt.Errorf("failed to update mission completion status: %s", err)
	}
	return nil
}

func (mdb *MissionsDb) List(offset, count int) ([]*Mission, error) {
	if offset < 0 {
		offset = 0
	}
	if count < 1 {
		count = 20
	}
	rows, err := mdb.Db.Query(`
		SELECT 
			id, assigned_cat_id, completion_status
		FROM
			missions
			LIMIT $1 OFFSET $2
	`, count, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of missions: %s", err)
	}
	defer rows.Close()

	missions := make([]*Mission, 0)
	for rows.Next() {
		var mission Mission
		var nullableCatId sql.NullString
		err = rows.Scan(&mission.Id, &nullableCatId, &mission.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mission row: %s", err)
		}
		if nullableCatId.Valid {
			mission.AssignedCatId = nullableCatId.String
		}
		missions = append(missions, &mission)
	}
	return missions, nil
}

func (mdb *MissionsDb) Get(uuid string) (*Mission, error) {
	var mission Mission
	var nullableCatId sql.NullString
	err := mdb.Db.QueryRow(`
		SELECT 
			id, assigned_cat_id, completion_status
		FROM missions 
		WHERE id = $1
	`, uuid).Scan(&mission.Id, &nullableCatId, &mission.Status)
	if nullableCatId.Valid {
		mission.AssignedCatId = nullableCatId.String
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mission: %s", err)
	}
	return &mission, nil
}

func (mdb *MissionsDb) Assign(uuid, catUuid string) error {
	_, err := mdb.Db.Exec(`
		UPDATE missions SET
			assigned_cat_id = $2
		WHERE id = $1
	`, uuid, catUuid)
	if err != nil {
		return fmt.Errorf("failed to update mission completion status: %s", err)
	}
	return nil
}
