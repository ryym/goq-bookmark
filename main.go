package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	entriesR := repo.NewEntriesRepo(db)

	e := echo.New()
	e.Renderer = NewTemplate("views/*.html")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})
	e.GET("/users", ShowUsers(usersR))
	e.GET("/users/:user_id", ShowBookmarks(usersR, bookmarksR, entriesR))
	e.POST("/users/:user_id", CreateBookmark(bookmarksR))
	e.GET("/bookmarks/:bookmark_id", EditBookmark(bookmarksR))
	e.POST("/bookmarks/:bookmark_id", UpdateBookmark(bookmarksR, usersR))
	e.GET("/entries", ShowEntries(entriesR))
	e.POST("/entries", CreateEntry(entriesR))
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
	entriesR *repo.EntriesRepo,
) echo.HandlerFunc {
	type bookmark struct {
		Bookmark model.Bookmark
		Entry    model.Entry
	}
	type data struct {
		User      model.User
		Bookmarks []bookmark
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

func CreateBookmark(bookmarksR *repo.BookmarksRepo) echo.HandlerFunc {
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

func EditBookmark(bookmarksR *repo.BookmarksRepo) echo.HandlerFunc {
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

func UpdateBookmark(bookmarksR *repo.BookmarksRepo, usersR *repo.UsersRepo) echo.HandlerFunc {
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

type entriesPageData struct {
	Entries []model.Entry
	Entry   model.Entry
	Errs    []error
}

func ShowEntries(entriesR *repo.EntriesRepo) echo.HandlerFunc {
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

func CreateEntry(entriesR *repo.EntriesRepo) echo.HandlerFunc {
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
