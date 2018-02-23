package repo

import (
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/goq"
)

type EntriesRepo struct {
	*Repo
}

func NewEntriesRepo(db *goq.DB) *EntriesRepo {
	return &EntriesRepo{newRepo(db)}
}

func (r *EntriesRepo) All() ([]model.Entry, error) {
	z := r.Builder
	q := z.Select(z.Entries.All()).From(z.Entries).OrderBy(z.Entries.ID.Desc())

	var entries []model.Entry
	err := r.DB.Query(q).Collect(z.ToSlice(&entries))
	return entries, err
}

func (r *EntriesRepo) Create(entry *model.Entry) error {
	z := r.Builder
	q := z.InsertInto(z.Entries, z.Entries.Except(z.Entries.ID).Columns()...).Values(entry)
	_, err := r.DB.Exec(q)
	return err
}
