package main

import (
	"github.com/go-martini/martini"
	"code.google.com/p/leveldb-go/leveldb"
	"log"
	"code.google.com/p/leveldb-go/leveldb/db"
)

func main() {
	m := martini.Classic()

	var opts db.Options
	level, err := leveldb.Open("hogeDB", &opts)
	if err != nil {
		panic(err)
	}
	defer level.Close()

	m.Map(level)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/register/:key/:value", RegisterData)
	m.Group("/read", func(r martini.Router) {
		r.Get("/:key", ReadData)
	})
	m.Run()
}

func RegisterData(params martini.Params, l *log.Logger, level *leveldb.DB) (int, string) {
	key := params["key"]
	value := params["value"]
	l.Println(key, value)

	var opts db.WriteOptions
	err := level.Set([]byte(key), []byte(value), &opts)
	if err != nil {
		message := "Error! " + err.Error() + " key = " + key + " value = " + value
		l.Fatal(message)
		return 400, message
	}
	return 200, "Hello!" + value
}

func ReadData(params martini.Params, l *log.Logger, level *leveldb.DB) (int, string) {
	key := params["key"]
	var opts db.ReadOptions
	value, err := level.Get([]byte(key), &opts)
	if err != nil {
		message := "Error! " + err.Error() + " key = " + key
		l.Fatal(message)
		return 400, message
	}

	return 200, "Data key = " + key + " value = " + string(value)
}
