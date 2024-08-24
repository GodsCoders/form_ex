package dbops

import (
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/rohanthewiz/serr"
)

func SavePerson(name string, age int) (err error) {
	db, err := sql.Open("duckdb", "coders.duck")
	if err != nil {
		fmt.Println(err)
		return
	}

	// We should *always* check for errors
	_, err = db.Exec(`CREATE TABLE if not exists person (age INTEGER, name VARCHAR)`)
	if err != nil {
		return serr.Wrap(err)
	}

	_, err = db.Exec(`INSERT INTO person (age, name) VALUES (?, ?)`, age, name)
	if err != nil {
		return serr.Wrap(err)
	}

	rs, err := db.Query(`SELECT age, name FROM person`)
	if err != nil {
		return serr.Wrap(err)
	}
	defer rs.Close()

	for rs.Next() {
		var dAge int
		var dName string

		err = rs.Scan(&dAge, &dName)
		if err != nil {
			return serr.Wrap(err)
		}

		fmt.Println("Received from DB -> age:", dAge, "name:", dName)
	}

	return
}
