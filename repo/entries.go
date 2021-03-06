package repo

import (
	"fmt"

	"github.com/ryym/goq"
	"github.com/ryym/goq-bookmark/model"
)

type EntriesRepo struct {
	*Repo
}

func NewEntriesRepo(db *goq.DB) *EntriesRepo {
	return &EntriesRepo{newRepo(db)}
}

func (r *EntriesRepo) Find(entryID int) (model.Entry, error) {
	q := z.Select(z.Entries.All()).From(z.Entries).Where(z.Entries.ID.Eq(entryID))
	var entry model.Entry
	err := r.DB.Query(q).First(z.Entries.ToElem(&entry))
	if err == nil && entry.ID == 0 {
		err = fmt.Errorf("could not find entry %d", entryID)
	}
	return entry, err
}

func (r *EntriesRepo) All() ([]model.Entry, error) {
	q := z.Select(z.Entries.All()).From(z.Entries).OrderBy(z.Entries.CreatedAt)

	var entries []model.Entry
	err := r.DB.Query(q).Collect(z.ToSlice(&entries))
	return entries, err
}

func (r *EntriesRepo) Create(entry *model.Entry) error {
	q := z.InsertInto(
		z.Entries,
		z.Entries.Except(z.Entries.ID, z.Entries.CreatedAt).Columns()...,
	).Values(entry)
	_, err := r.DB.Exec(q)
	return err
}

func (r *EntriesRepo) Update(entry *model.Entry) error {
	q := z.Update(z.Entries).Elem(entry, z.Entries.Title, z.Entries.URL)
	_, err := r.DB.Exec(q)
	return err
}

func (r *EntriesRepo) Delete(entryID int) error {
	q := z.DeleteFrom(z.Entries).Where(z.Entries.ID.Eq(entryID))
	_, err := r.DB.Exec(q)
	return err
}
