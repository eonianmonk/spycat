package main

import (
	"database/sql"
	"os"

	"github.com/eonianmonk/spycat"
	"github.com/eonianmonk/spycat/internal/data"
	"github.com/eonianmonk/spycat/internal/http"
	"github.com/eonianmonk/spycat/internal/http/context"

	_ "github.com/lib/pq"
)

func main() {
	db := connDb()
	validator, err := spycat.NewCatValidator()
	if err != nil {
		panic(err)
	}
	http.Run(&context.DbsCtx{
		CatsDb:     &data.CatsDb{Db: db},
		MissionsDb: &data.MissionsDb{Db: db},
		TargetsDb:  &data.TargetDb{Db: db},
	}, validator, port)
}

func connDb() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv(dbConnStr))
	if err != nil {
		panic(err)
	}
	return db
}
