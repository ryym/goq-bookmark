package main

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/ryym/goq"
	"github.com/ryym/goq-bookmark/ctx"
	"github.com/ryym/goq-bookmark/handler"
	"github.com/ryym/goq-bookmark/repo"
)

func main() {
	db, err := connectToDB()
	if err != nil {
		panic(err)
	}

	repo.Init(db)
	app := ctx.NewAppContext(db)

	e := echo.New()
	defineRoutes(e, app)
	e.Renderer = NewTemplate("views/*.html")
	e.Logger.Fatal(e.Start(":8000"))
}

func connectToDB() (*goq.DB, error) {
	conn := "port=5431 user=bookmark sslmode=disable"
	db, err := goq.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	nTries := 3
	for i := 0; i < nTries; i++ {
		if err = db.DB.Ping(); err == nil {
			break
		}
		fmt.Println("waiting for DB...")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %s", err)
	}

	return db, nil
}

func defineRoutes(e *echo.Echo, app *ctx.AppContext) {
	e.GET("/", handler.Home())
	e.GET("/users", handler.ShowUsers(app))
	e.GET("/users/:user_id", handler.ShowBookmarks(app))
	e.POST("/users/:user_id", handler.CreateBookmark(app))
	e.GET("/bookmarks/:bookmark_id", handler.EditBookmark(app))
	e.POST("/bookmarks/:bookmark_id", handler.UpdateBookmark(app))
	// XXX: You should not use GET for deletion.
	e.GET("/bookmarks/:bookmark_id/delete", handler.DeleteBookmark(app))
	e.GET("/entries", handler.ShowEntries(app))
	e.POST("/entries", handler.CreateEntry(app))
	e.GET("/entries/:entry_id", handler.EditEntry(app))
	e.POST("/entries/:entry_id", handler.UpdateEntry(app))
	// XXX: You should not use GET for deletion.
	e.GET("/entries/:entry_id/delete", handler.DeleteEntry(app))
}
