package repo

import (
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/goq"
)

type BookmarksRepo struct {
	DB      *goq.DB
	Builder *Builder
}

func NewBookmarksRepo(db *goq.DB) *BookmarksRepo {
	return &BookmarksRepo{db, NewBuilder(db.Dialect())}
}

func (r *BookmarksRepo) FromUser(userID int) ([]model.Bookmark, []model.Entry, error) {
	z := r.Builder
	q := z.Select(
		z.Bookmarks.All(),
		z.Entries.All(),
	).From(
		z.Bookmarks,
	).Joins(
		z.Bookmarks.Entries(z.Entries),
	).Where(
		z.Bookmarks.UserID.Eq(userID),
	)

	var bookmarks []model.Bookmark
	var entries []model.Entry
	err := r.DB.Query(q).Collect(
		z.Bookmarks.ToSlice(&bookmarks),
		z.Entries.ToSlice(&entries),
	)
	return bookmarks, entries, err
}
