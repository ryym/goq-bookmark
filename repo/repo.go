package repo

import "github.com/ryym/goq"

var z *Builder

type Repo struct {
	DB *goq.DB
}

func newRepo(db *goq.DB) *Repo {
	if z == nil {
		panic("Query builder does not be initialized")
	}
	return &Repo{db}
}

func Init(db *goq.DB) {
	z = NewBuilder(db.Dialect())
}
