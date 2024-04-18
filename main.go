package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"net/http"
)

const DBNAME = "./foo.db"

func main() {
  db := setupDatabase()
  defer db.Close()

  command := os.Args[1]
  if command == "server" {
    runServer(db)
  } else {
    myType := checkType(command)
    insertBy(db, command, myType)
    dumpDatabase(db)
  }
}

func runServer(db *sql.DB) {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    myType := checkType(url)
    insertBy(db, url, myType)
    dumpDatabase(db)
    fmt.Fprintf(w, "%q", url)
  })

  log.Fatal(http.ListenAndServe(":6969", nil))
}

func checkType(input string) string {
	switch {
	case strings.Contains(input, "youtube") || strings.Contains(input, "youtu.be"):
		return "youtube"
	default:
		log.Printf("Couldn't match type for `%s`", input)
		return "unknown"
	}
}

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", DBNAME)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
  CREATE TABLE IF NOT EXISTS foo(
    id integer not null primary key autoincrement, 
    name text,
    type text
  );
  `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}

	return db
}

func insertBy(db *sql.DB, name, t_type string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into foo(name, type) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, t_type)
	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}
}

func dumpDatabase(db *sql.DB) {
	rows, err := db.Query("select id, name, type from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var type_t string

		err = rows.Scan(&id, &name, &type_t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, type_t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
