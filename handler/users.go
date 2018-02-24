package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/go-bookmark/repo"
)

func ShowUsers(app interface {
	UsersRepo() *repo.UsersRepo
}) echo.HandlerFunc {
	usersR := app.UsersRepo()
	return func(c echo.Context) error {
		users, err := usersR.All()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "users", users)
	}
}

func ShowBookmarks(app interface {
	UsersRepo() *repo.UsersRepo
	BookmarksRepo() *repo.BookmarksRepo
}) echo.HandlerFunc {
	type bookmark struct {
		Bookmark model.Bookmark
		Entry    model.Entry
	}
	type data struct {
		User      model.User
		Bookmarks []bookmark
		Entries   []model.Entry
	}

	usersR := app.UsersRepo()
	bookmarksR := app.BookmarksRepo()
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

		bks := make([]bookmark, 0, len(bookmarks))
		for i, bk := range bookmarks {
			bks = append(bks, bookmark{bk, entries[i]})
		}

		candidates, err := bookmarksR.UnbookmarkedEntries(userID)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "bookmarks", &data{
			User:      user,
			Bookmarks: bks,
			Entries:   candidates,
		})
	}
}

func CreateBookmark(app interface {
	BookmarksRepo() *repo.BookmarksRepo
}) echo.HandlerFunc {
	bookmarksR := app.BookmarksRepo()
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			return fmt.Errorf("Invalid user ID: %s", err)
		}
		entryID, err := strconv.Atoi(c.FormValue("entry_id"))
		if err != nil {
			return fmt.Errorf("Invalid entry ID: %s", err)
		}

		bookmark := model.Bookmark{
			UserID:  userID,
			EntryID: entryID,
			Comment: strings.TrimSpace(c.FormValue("comment")),
		}

		err = bookmarksR.Create(&bookmark)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d", userID))
	}
}
