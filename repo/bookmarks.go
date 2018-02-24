package repo

import (
	"fmt"

	"github.com/ryym/goq-bookmark/model"
	"github.com/ryym/goq"
)

type BookmarksRepo struct {
	*Repo
}

func NewBookmarksRepo(db *goq.DB) *BookmarksRepo {
	return &BookmarksRepo{newRepo(db)}
}

func (r *BookmarksRepo) Find(bookmarkID int) (model.Bookmark, error) {
	q := z.Select(z.Bookmarks.All()).From(z.Bookmarks).Where(z.Bookmarks.ID.Eq(bookmarkID))

	var bookmark model.Bookmark
	err := r.DB.Query(q).First(z.Bookmarks.ToElem(&bookmark))
	if err == nil && bookmark.ID == 0 {
		err = fmt.Errorf("Could not find bookmark %d", bookmarkID)
	}

	return bookmark, err
}

func (r *BookmarksRepo) FindWithAssocs(bookmarkID int) (model.Bookmark, model.User, model.Entry, error) {
	b, u, e := z.Bookmarks.As("b"), z.Users.As("u"), z.Entries.As("e")
	q := z.Select(b.All(), u.All(), e.All()).From(b).Joins(
		b.Users(u),
		b.Entries(e),
	).Where(b.ID.Eq(bookmarkID))

	var bookmark model.Bookmark
	var user model.User
	var entry model.Entry
	err := r.DB.Query(q).First(
		b.ToElem(&bookmark),
		u.ToElem(&user),
		e.ToElem(&entry),
	)

	if err == nil && bookmark.ID == 0 {
		err = fmt.Errorf("Could not find bookmark %d", bookmarkID)
	}

	return bookmark, user, entry, err
}

func (r *BookmarksRepo) FromUser(userID int) ([]model.Bookmark, []model.Entry, error) {
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
	q := z.Select(z.Entries.All()).From(z.Entries).Where(
		z.Entries.ID.NotIn(
			z.Select(z.Bookmarks.EntryID).From(z.Bookmarks).Where(
				z.Bookmarks.UserID.Eq(userID),
			),
		),
	)

	var entries []model.Entry
	err := r.DB.Query(q).Collect(z.Entries.ToSlice(&entries))
	return entries, err
}

func (r *BookmarksRepo) Create(bookmark *model.Bookmark) error {
	q := z.InsertInto(
		z.Bookmarks,
		z.Bookmarks.Except(z.Bookmarks.ID).Columns()...,
	).Values(bookmark)
	_, err := r.DB.Exec(q)
	return err
}

func (r *BookmarksRepo) Update(bookmark *model.Bookmark) error {
	q := z.Update(z.Bookmarks).Elem(bookmark, z.Bookmarks.Comment)
	_, err := r.DB.Exec(q)
	return err
}

func (r *BookmarksRepo) Delete(bookmarkID int) error {
	q := z.DeleteFrom(z.Bookmarks).Where(z.Bookmarks.ID.Eq(bookmarkID))
	_, err := r.DB.Exec(q)
	return err
}
