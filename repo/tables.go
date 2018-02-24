package repo

import "github.com/ryym/goq-bookmark/model"

//go:generate goq

type Tables struct {
	users     model.User
	entries   model.Entry
	bookmarks model.Bookmark
}
