package data

import (
	"database/sql"
	"fmt"

	"github.com/eonianmonk/spycat"

	"github.com/shopspring/decimal"
)

type Cat struct {
	Id                string          `json:"id,omitempty"`
	Name              string          `json:"name,omitempty"`
	YearsOfExperience int             `json:"years_of_experience,omitempty"`
	Breed             spycat.Breed    `json:"breed,omitempty"`
	Salary            decimal.Decimal `json:"salary,omitempty"`
}

type CatsDb struct {
	Db *sql.DB
}

// creates new cat, returns cat data with generated id
// cat's breed should be validated externally
func (catsDb *CatsDb) Create(cat *Cat) (*Cat, error) {
	rows, err := catsDb.Db.Query(`
		insert into cats (
			name, years_of_experience, breed, salary
		) values ($1,$2,$3,$4) returning id`,
		cat.Name, cat.YearsOfExperience, string(cat.Breed), cat.Salary.String())
	if err != nil {
		return nil, fmt.Errorf("failed to insert cat row: %s", err)
	}

	var uuid string

	for rows.Next() {
		if err := rows.Scan(&uuid); err != nil {
			return nil, fmt.Errorf("failed to scan registered cat uuid: %s", err)
		}
	}
	cat.Id = uuid
	return cat, nil
}

// deletes cat row by cat's uuid
func (catsDb *CatsDb) Delete(uuid string) error {
	_, err := catsDb.Db.Exec("delete from cats where id = $1", uuid)
	if err != nil {
		return fmt.Errorf("failed to delete cat by id(%s): %s", uuid, err)
	}
	return nil
}

// update cats salary
func (catsDb *CatsDb) UpdateSalary(uuid string, salary decimal.Decimal) error {
	_, err := catsDb.Db.Exec(`
		UPDATE cats SET
			salary = $1
		WHERE id = $2
	`, salary.String(), uuid)
	if err != nil {
		return fmt.Errorf("failed to update salary for %s: %s", uuid, err)
	}
	return nil
}

// list cats
func (cdb *CatsDb) List(offset int, count int) ([]*Cat, error) {
	if offset < 0 {
		offset = 0
	}
	if count <= 0 {
		count = 20
	}
	rows, err := cdb.Db.Query(`
		SELECT 
			id, name, years_of_experience, breed, salary
		FROM
			cats
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`, offset, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of cats: %s", err)
	}
	defer rows.Close()

	cats := make([]*Cat, 0)
	for rows.Next() {
		var cat Cat
		var salaryString string
		err := rows.Scan(&cat.Id, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &salaryString)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cat row: %s", err)
		}

		cat.Salary, err = decimal.NewFromString(salaryString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse salary: %s", err)
		}
		cats = append(cats, &cat)
	}
	return cats, nil
}

// gets single cat by id
func (cdb *CatsDb) GetCat(uuid string) (*Cat, error) {
	var cat Cat
	var salaryString string

	err := cdb.Db.QueryRow(`
		Select 
			id, name, years_of_experience, breed, salary
		from
			cats
		where
			id = $1
	`, uuid).Scan(&cat.Id, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &salaryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get cat %s: %s", uuid, err)
	}
	cat.Salary, err = decimal.NewFromString(salaryString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse salary: %s", err)
	}
	return &cat, nil
}
