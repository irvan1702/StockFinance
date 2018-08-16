package config

import (
	"database/sql"
	"fmt"
	"html/template"

	_ "github.com/lib/pq"
)

var Db *sql.DB
var Tpl *template.Template

func Init() {
	var err error
	Db, err = sql.Open("postgres", "postgres://postgres:password@/StockFinance?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = Db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Database Connected")
}
