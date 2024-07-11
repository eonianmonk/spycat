package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/eonianmonk/spycat"
)

type Target struct {
	Id        string                 `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Country   string                 `json:"country,omitempty"`
	Status    spycat.ComletionStatus `json:"status,omitempty"`
	Notes     string                 `json:"notes,omitempty"`
	MissionId string                 `json:"mission_id"`
}

type TargetDb struct {
	Db *sql.DB
}

func buildBulkTargetInsertValues(targets []*Target) (string, interface{}) {
	b := strings.Builder{}
	// because target has 6 fields (one of which is id and currently unknown)
	args := make([]interface{}, 0, 5*len(targets))
	for i := range targets {
		b.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5))
		args = append(args, targets[i].Name, targets[i].Country, targets[i].Status, targets[i].Notes, targets[i].MissionId)
	}
	valuesQ := b.String()
	// remove last ','
	valuesQ = valuesQ[:len(valuesQ)-1]
	return valuesQ, args
}

// creates targets via provided txs
func (tdb *TargetDb) CreateMany(targets []*Target, missionId string, tx *sql.Tx) ([]*Target, error) {

	localTx := false
	var err error
	if tx == nil {
		tx, err = tdb.Db.Begin()
		if err != nil {
			return nil, err
		}
		localTx = true
	}

	if missionId != "" {
		for i := range len(targets) {
			targets[i].MissionId = missionId
		}
	}
	valQ, args := buildBulkTargetInsertValues(targets)

	rows, err := tx.Query(fmt.Sprintf(`
		INSERT INTO targets (
			name, country, completion_statust, notes,mission_id
		) VALUES %s 
		RETURNING id,name, country, completion_statust, notes,mission_id;
	`, valQ), args)
	if err != nil {
		return nil, fmt.Errorf("failed to insert targets: %s", err)
	}
	defer rows.Close()

	targetsWithId := make([]*Target, len(targets))
	i := 0
	for rows.Next() {
		var target Target
		err := rows.Scan(&target.Id, &target.Name, &target.Country,
			&target.Status, &target.Notes, &target.MissionId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan target row: %s", err)
		}
		targetsWithId[i] = &target
		i++
	}

	if localTx {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return targetsWithId, nil
}

func (tdb *TargetDb) Delete(uuid string) error {
	_, err := tdb.Db.Query(`DELETE FROM targets where id = $1`, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete target %s: %s", uuid, err)
	}
	return nil
}

func buildBulkValuesFromMissions(missions []*Mission) (string, []interface{}) {
	vals := make([]string, 0, len(missions))
	args := make([]interface{}, 0, len(missions))
	ix := 1
	for _, mission := range missions {
		vals = append(vals, fmt.Sprintf("$%d", ix))
		args = append(args, mission.Id)
		ix++
	}
	return fmt.Sprintf("(%s)", strings.Join(vals, ",")), args
}

func (tdb *TargetDb) GetTargetsForMissions(missions []*Mission) error {
	vals, args := buildBulkValuesFromMissions(missions)
	rows, err := tdb.Db.Query(fmt.Sprintf(`
		SELECT 
			id, name, country, completion_status, notes, mission_id
		FROM 
			targets
		WHERE id IN %s 
	`, vals), args)
	if err != nil {
		return fmt.Errorf("failed to get targets for missions: %s", err)
	}
	defer rows.Close()

	// map of mission ids to it's targetss
	missionsToTargets := make(map[string][]*Target)

	for rows.Next() {
		var target Target

		err = rows.Scan(&target.Id, &target.Name, &target.Country, &target.Status, &target.Notes, &target.MissionId)
		if err != nil {
			return fmt.Errorf("failed to scan targets row: %s", err)
		}
		_, ok := missionsToTargets[target.MissionId]
		if !ok {
			missionsToTargets[target.MissionId] = make([]*Target, 0, 1)
		}
		missionsToTargets[target.MissionId] = append(missionsToTargets[target.MissionId], &target)
	}

	// assign targets to respective missions
	for i := range missions {
		_, ok := missionsToTargets[missions[i].Id]
		// reverse is possible and is normal behaviour (mission without targets)
		if ok {
			missions[i].Targets = missionsToTargets[missions[i].Id]
		}
	}
	return nil
}

func (tdb *TargetDb) Update(uuid, note string, status spycat.ComletionStatus) error {
	_, err := tdb.Db.Exec(`
		UPDATE targets SET
			completion_status = $2,
			notes = $3
		WHERE id = $1
	`, uuid, status, note)
	if err != nil {
		return fmt.Errorf("failed to update target state: %s", err)
	}
	return nil
}
