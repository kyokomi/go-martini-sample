package main

import (
	"github.com/go-martini/martini"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"fmt"
)

func main() {
	m := martini.Classic()

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `
	create table IF NOT EXISTS foo (id integer not null primary key, name text);
	`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
		return
	}

	m.Map(db)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/register/:key/:value", RegisterData)
	m.Group("/read", func(r martini.Router) {
		r.Get("/:key", ReadData)
	})
	m.Run()
}

func RegisterData(params martini.Params, l *log.Logger, db *sql.DB) (int, string) {
	key := params["key"]
	value := params["value"]
	l.Println(key, value)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("[%s = %s] ねぷねぷのぷりん%03d", key, value, i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	return 200, "Hello!" + value
}

func ReadData(params martini.Params, l *log.Logger, db *sql.DB) (int, string) {
	key := params["key"]

	stmt, err := db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow(key).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	return 200, "Data key = " + key + " name = " + name
}
