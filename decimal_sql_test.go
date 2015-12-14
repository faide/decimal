package decimal

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func TestSQLiteFloat(t *testing.T) {
	var (
		Db     *sqlx.DB
		id1    = 1
		id2    = 2
		amnt1  = NewFromFloat(55.33)
		amnt2  = New(0, 0)
		schema = `
			CREATE TABLE IF NOT EXISTS some_record
				(
				id integer,
				amount numeric
				);`
	)

	Db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln(err)
	}
	err = Db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	// create the schema
	tx := Db.MustBegin()
	tx.MustExec(schema)
	tx.Commit()

	// set a value
	tx = Db.MustBegin()
	tx.Exec("INSERT INTO some_record (id, amount) VALUES ($1, $2)", id1, amnt1)
	tx.Exec("INSERT INTO some_record (id, amount) VALUES ($1, $2)", id2, amnt2)
	if err = tx.Commit(); err != nil {
		t.Fatalf("INSERT commit failed")
	}

	// use previously set value to query and test our scanner
	amount := Decimal{}
	err = Db.QueryRow(
		"SELECT amount FROM some_record WHERE id==?1", id1).Scan(&amount)

	switch err {

	case sql.ErrNoRows:
		t.Fatalf("Expected to find a row..., but none found")

	case nil:
		if !amount.Equals(amnt1) {
			t.Fatalf("Expected to find 54.33 got %s", amount)
		}

	default:
		t.Fatalf("Expected to find a row with a float but got err", err)

	}

	err = Db.QueryRow(
		"SELECT amount FROM some_record WHERE id==?1", id2).Scan(&amount)

	switch err {

	case sql.ErrNoRows:
		t.Fatalf("Expected to find a row..., but none found")

	case nil:
		if !amount.Equals(amnt2) {
			t.Fatalf("Expected to find zero got", amount)
		}

	default:
		t.Fatalf("Expected to find a row with a float but got err", err)
	}
}
