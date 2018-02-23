package main

import (
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/ryym/go-bookmark/repo"
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

	usersR := repo.NewUsersRepo(db)

	e := echo.New()
	e.Renderer = NewTemplate("views/*.html")
	e.GET("/", ShowUsers(usersR))
	e.Logger.Fatal(e.Start(":8000"))
}

func ShowUsers(usersR *repo.UsersRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := usersR.All()
		if err != nil {
			panic(err)
		}
		return c.Render(http.StatusOK, "users", users)
	}
}
