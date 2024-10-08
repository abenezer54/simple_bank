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
	dbSource = "postgresql://root:root123@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testQueriesTx *Store

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)
	testQueriesTx = NewStore(conn)

	exitCode := m.Run()
	os.Exit(exitCode)
}
