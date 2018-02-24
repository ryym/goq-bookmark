package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/ryym/goq-bookmark/model"
	"github.com/ryym/goq-bookmark/repo"
)

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
