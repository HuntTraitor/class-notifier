package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func newTestDB(t *testing.T) *sql.DB {

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(".env file could not be loaded")
	}
	dsn := os.Getenv("TEST_POSTGRES_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	//calls when function goes out of scope
	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}
