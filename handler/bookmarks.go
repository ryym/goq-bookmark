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

type editBookmarkPageData struct {
	Bookmark model.Bookmark
	User     model.User
	Entry    model.Entry
}

func EditBookmark(app interface {
	BookmarksRepo() *repo.BookmarksRepo
}) echo.HandlerFunc {
	bookmarksR := app.BookmarksRepo()
	return func(c echo.Context) error {
		bookmarkID, err := strconv.Atoi(c.Param("bookmark_id"))
		if err != nil {
			return fmt.Errorf("Invalid bookmark ID: %s", err)
		}

		bookmark, user, entry, err := bookmarksR.FindWithAssocs(bookmarkID)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "bookmark-edit", &editBookmarkPageData{
			Bookmark: bookmark,
			Entry:    entry,
			User:     user,
		})
	}
}

func UpdateBookmark(app interface {
	BookmarksRepo() *repo.BookmarksRepo
}) echo.HandlerFunc {
	bookmarksR := app.BookmarksRepo()
	return func(c echo.Context) error {
		bookmarkID, err := strconv.Atoi(c.Param("bookmark_id"))
		if err != nil {
			return fmt.Errorf("Invalid bookmark ID: %s", err)
		}

		bookmark, err := bookmarksR.Find(bookmarkID)
		if err != nil {
			return err
		}

		bookmark.Comment = strings.TrimSpace(c.FormValue("comment"))
		err = bookmarksR.Update(&bookmark)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d", bookmark.UserID))
	}
}

func DeleteBookmark(app interface {
	BookmarksRepo() *repo.BookmarksRepo
}) echo.HandlerFunc {
	bookmarksR := app.BookmarksRepo()
	return func(c echo.Context) error {
		bookmarkID, err := strconv.Atoi(c.Param("bookmark_id"))
		if err != nil {
			return fmt.Errorf("Invalid bookmark ID: %s", err)
		}
		err = bookmarksR.Delete(bookmarkID)
		if err != nil {
			return err
		}

		referer := c.Request().Header.Get("referer")
		if referer == "" {
			referer = "/"
		}
		return c.Redirect(http.StatusFound, referer)
	}
}
