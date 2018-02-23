package repo

import "github.com/ryym/goq"

type Repo struct {
	DB      *goq.DB
	Builder *Builder
}

func newRepo(db *goq.DB) *Repo {
	return &Repo{db, NewBuilder(db.Dialect())}
}
