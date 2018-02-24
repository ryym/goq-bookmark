package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/ryym/go-bookmark/ctx"
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

	appCtx := ctx.NewAppContext(db)

	e := echo.New()
	e.Renderer = NewTemplate("views/*.html")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})
	e.GET("/users", ShowUsers(appCtx))
	e.GET("/users/:user_id", ShowBookmarks(appCtx))
	e.POST("/users/:user_id", CreateBookmark(appCtx))
	e.GET("/bookmarks/:bookmark_id", EditBookmark(appCtx))
	e.POST("/bookmarks/:bookmark_id", UpdateBookmark(appCtx))

	// XXX: You should not use GET for deletion.
	e.GET("/bookmarks/:bookmark_id/delete", DeleteBookmark(appCtx))

	e.GET("/entries", ShowEntries(appCtx))
	e.POST("/entries", CreateEntry(appCtx))
	e.GET("/entries/:entry_id", EditEntry(appCtx))
	e.POST("/entries/:entry_id", UpdateEntry(appCtx))

	// XXX: You should not use GET for deletion.
	e.GET("/entries/:entry_id/delete", DeleteEntry(appCtx))

	e.Logger.Fatal(e.Start(":8000"))
}

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

type entriesPageData struct {
	Entries []model.Entry
	Entry   model.Entry
	Errs    []error
}

func ShowEntries(app interface {
	EntriesRepo() *repo.EntriesRepo
}) echo.HandlerFunc {
	entriesR := app.EntriesRepo()
	return func(c echo.Context) error {
		entries, err := entriesR.All()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "entries", &entriesPageData{
			Entries: entries,
		})
	}
}

func CreateEntry(app interface {
	EntriesRepo() *repo.EntriesRepo
}) echo.HandlerFunc {
	entriesR := app.EntriesRepo()
	return func(c echo.Context) error {
		entry := model.Entry{
			Title: c.FormValue("title"),
			URL:   c.FormValue("url"),
		}

		errs := model.ValidateEntry(&entry)
		if len(errs) > 0 {
			entries, err := entriesR.All()
			if err != nil {
				return err
			}
			return c.Render(http.StatusOK, "entries", &entriesPageData{
				Entries: entries,
				Entry:   entry,
				Errs:    errs,
			})
		}

		err := entriesR.Create(&entry)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/entries")
	}
}

type editEntryPageData struct {
	Errs  []error
	Entry model.Entry
}

func EditEntry(app interface {
	EntriesRepo() *repo.EntriesRepo
}) echo.HandlerFunc {
	entriesR := app.EntriesRepo()
	return func(c echo.Context) error {
		entryID, err := strconv.Atoi(c.Param("entry_id"))
		if err != nil {
			return fmt.Errorf("Invalid entry ID: %s", err)
		}

		entry, err := entriesR.Find(entryID)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "entry-edit", &editEntryPageData{
			Entry: entry,
		})
	}
}

func UpdateEntry(app interface {
	EntriesRepo() *repo.EntriesRepo
}) echo.HandlerFunc {
	entriesR := app.EntriesRepo()
	return func(c echo.Context) error {
		entryID, err := strconv.Atoi(c.Param("entry_id"))
		if err != nil {
			return fmt.Errorf("Invalid entry ID: %s", err)
		}

		entry := model.Entry{
			ID:    entryID,
			Title: c.FormValue("title"),
			URL:   c.FormValue("url"),
		}
		errs := model.ValidateEntry(&entry)
		if len(errs) > 0 {
			return c.Render(http.StatusOK, "entry-edit", &editEntryPageData{
				Entry: entry,
				Errs:  errs,
			})
		}

		err = entriesR.Update(&entry)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, "/entries")
	}
}

func DeleteEntry(app interface {
	EntriesRepo() *repo.EntriesRepo
}) echo.HandlerFunc {
	entriesR := app.EntriesRepo()
	return func(c echo.Context) error {
		entryID, err := strconv.Atoi(c.Param("entry_id"))
		if err != nil {
			return fmt.Errorf("Invalid entry ID: %s", err)
		}

		err = entriesR.Delete(entryID)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, "/entries")
	}
}
