package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/go-bookmark/repo"
	"github.com/ryym/goq"
)

func main() {
	conn := "port=5431 user=bookmark sslmode=disable"
	db, err := goq.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	usersR := repo.NewUsersRepo(db)
	bookmarksR := repo.NewBookmarksRepo(db)

	e := echo.New()
	e.Renderer = NewTemplate("views/*.html")
	e.GET("/", ShowUsers(usersR))
	e.GET("/:user_id", ShowBookmarks(usersR, bookmarksR))
	e.Logger.Fatal(e.Start(":8000"))
}

func ShowUsers(usersR *repo.UsersRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := usersR.All()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "users", users)
	}
}

func ShowBookmarks(
	usersR *repo.UsersRepo,
	bookmarksR *repo.BookmarksRepo,
) echo.HandlerFunc {
	type data struct {
		User      model.User
		Bookmarks []model.Bookmark
		Entries   []model.Entry
	}

	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			return err
		}

		user, err := usersR.Find(userID)
		if err != nil {
			return err
		}
		if user.ID == 0 {
			return fmt.Errorf("invalid user ID: %d", userID)
		}

		bookmarks, entries, err := bookmarksR.FromUser(userID)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "bookmarks", &data{
			User:      user,
			Bookmarks: bookmarks,
			Entries:   entries,
		})
	}
}
