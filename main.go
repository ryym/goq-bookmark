package main

import (
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/ryym/goq"
)

func main() {
	conn := "port=5432 database=bookmark sslmode=disable"
	db, err := goq.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	e := echo.New()
	e.Renderer = NewTemplate("views/*.html")
	e.Logger.Fatal(e.Start(":8000"))
}
