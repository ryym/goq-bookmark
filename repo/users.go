package repo

import (
	"fmt"

	"github.com/ryym/goq"
	"github.com/ryym/goq-bookmark/model"
)

type UsersRepo struct {
	*Repo
}

func NewUsersRepo(db *goq.DB) *UsersRepo {
	return &UsersRepo{newRepo(db)}
}

func (r *UsersRepo) Find(id int) (model.User, error) {
	q := z.Select(z.Users.All()).From(z.Users).Where(z.Users.ID.Eq(id))
	var user model.User
	err := r.DB.Query(q).First(z.Users.ToElem(&user))
	if err == nil && user.ID == 0 {
		err = fmt.Errorf("could not find user %d", id)
	}
	return user, err
}

func (r *UsersRepo) All() ([]model.User, error) {
	q := z.Select(z.Users.All()).From(z.Users).OrderBy(z.Users.ID)

	var users []model.User
	err := r.DB.Query(q).Collect(z.Users.ToSlice(&users))
	return users, err
}
