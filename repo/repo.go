package repo

import "github.com/ryym/goq"

// This sample application choose to put a query builder
// as a global variable for simplicity. But of course
// this is not a mandatory. You can declare a field on
// repository structs for a query builder and use it as well.
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
