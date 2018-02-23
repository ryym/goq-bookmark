package repo

import (
	"github.com/ryym/go-bookmark/model"
	"github.com/ryym/goq"
)

type UsersRepo struct {
	*Repo
}

func NewUsersRepo(db *goq.DB) *UsersRepo {
	return &UsersRepo{newRepo(db)}
}

func (r *UsersRepo) Find(id int) (model.User, error) {
	z := r.Builder
	q := z.Select(z.Users.All()).From(z.Users).Where(z.Users.ID.Eq(id))
	var user model.User
	err := r.DB.Query(q).First(z.Users.ToElem(&user))
	return user, err
}

func (r *UsersRepo) All() ([]model.User, error) {
	z := r.Builder
	q := z.Select(z.Users.All()).From(z.Users).OrderBy(z.Users.ID)

	var users []model.User
	err := r.DB.Query(q).Collect(z.Users.ToSlice(&users))
	return users, err
}
