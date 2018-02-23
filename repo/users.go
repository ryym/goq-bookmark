package repo

import (
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/goq"
)

type UsersRepo struct {
	DB      *goq.DB
	Builder *Builder
}

func NewUsersRepo(db *goq.DB) *UsersRepo {
	return &UsersRepo{db, NewBuilder(db.Dialect())}
}

func (r *UsersRepo) All() ([]model.User, error) {
	z := r.Builder
	q := z.Select(z.Users.All()).From(z.Users).OrderBy(z.Users.ID)

	var users []model.User
	err := r.DB.Query(q).Collect(z.Users.ToSlice(&users))
	return users, err
}
