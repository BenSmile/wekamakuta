package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSrouce = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSrouce)
	if err != nil {
		log.Fatal("we cannot connect to the db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())

}
