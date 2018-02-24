package ctx

import (
	"github.com/ryym/go-bookmark/repo"
	"github.com/ryym/goq"
)

// We don't use interfaces for repositories because
// this is just a sample app.

type AppContext struct {
	usersRepo     *repo.UsersRepo
	entriesRepo   *repo.EntriesRepo
	bookmarksRepo *repo.BookmarksRepo
}

func NewAppContext(db *goq.DB) *AppContext {
	return &AppContext{
		usersRepo:     repo.NewUsersRepo(db),
		bookmarksRepo: repo.NewBookmarksRepo(db),
		entriesRepo:   repo.NewEntriesRepo(db),
	}
}

func (c *AppContext) UsersRepo() *repo.UsersRepo {
	return c.usersRepo
}

func (c *AppContext) EntriesRepo() *repo.EntriesRepo {
	return c.entriesRepo
}

func (c *AppContext) BookmarksRepo() *repo.BookmarksRepo {
	return c.bookmarksRepo
}
