package repo

import (
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/goq"
)

type BookmarksRepo struct {
	*Repo
}

func NewBookmarksRepo(db *goq.DB) *BookmarksRepo {
	return &BookmarksRepo{newRepo(db)}
}

func (r *BookmarksRepo) FromUser(userID int) ([]model.Bookmark, []model.Entry, error) {
	z := r.Builder
	q := z.Select(z.Bookmarks.All(), z.Entries.All()).
		From(z.Bookmarks).
		Joins(z.Bookmarks.Entries(z.Entries)).
		Where(z.Bookmarks.UserID.Eq(userID))

	var bookmarks []model.Bookmark
	var entries []model.Entry
	err := r.DB.Query(q).Collect(
		z.Bookmarks.ToSlice(&bookmarks),
		z.Entries.ToSlice(&entries),
	)
	return bookmarks, entries, err
}

func (r *BookmarksRepo) UnbookmarkedEntries(userID int) ([]model.Entry, error) {
	z := r.Builder
	q := z.Select(z.Entries.All()).From(z.Entries).Where(z.Entries.ID.NotIn(
		z.Select(z.Bookmarks.EntryID).From(z.Bookmarks).Where(
			z.Bookmarks.UserID.Eq(userID),
		),
	))

	var entries []model.Entry
	err := r.DB.Query(q).Collect(z.Entries.ToSlice(&entries))
	return entries, err
}

func (r *BookmarksRepo) Create(bookmark *model.Bookmark) error {
	z := r.Builder
	q := z.InsertInto(
		z.Bookmarks,
		z.Bookmarks.Except(z.Bookmarks.ID).Columns()...,
	).Values(bookmark)
	_, err := r.DB.Exec(q)
	return err
}
