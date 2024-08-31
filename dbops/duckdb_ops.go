package dbops

import (
	"database/sql"
	"fmt"
	"form_ex/person"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/rohanthewiz/logger"
	"github.com/rohanthewiz/serr"
)

func SavePerson(name string, age int) (err error) {
	db, err := GetDBHandle()
	if err != nil {
		return serr.Wrap(err)
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

	return
}

// GetPersons retrieves a list of persons from the database.
func GetPersons() (persons []person.Person, err error) {
	db, err := GetDBHandle()
	if err != nil {
		return persons, serr.Wrap(err)
	}

	rs, err := db.Query(`SELECT age, name FROM person`)
	if err != nil {
		return persons, serr.Wrap(err)
	}
	defer rs.Close()

	for rs.Next() {
		var dAge int
		var dName string

		err = rs.Scan(&dAge, &dName)
		if err != nil {
			return persons, serr.Wrap(err)
		}

		p := person.Person{
			Age:  dAge,
			Name: dName,
		}

		persons = append(persons, p)
		fmt.Println("Received from DB -> age:", dAge, "name:", dName)
	}

	return
}

func GetDBHandle() (db *sql.DB, err error) {
	db, err = sql.Open("duckdb", "coders.duck")
	if err != nil {
		logger.LogErr(serr.Wrap(err))
		return
	}
	return db, nil
}

// *
// 		const htmlBegin = `<body style="background-color: slategrey;font-weight: bold;">
//			<table border=1 cellpadding=3>`
//		const htmlEnd = `</table></body>`
//
//		row := fmt.Sprintf(`<tr><td style="color:#3adeda">Name: %s</td><td> Age: %d</td></tr>`,
//			name, age)
// */
