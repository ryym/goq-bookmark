package repo

import "github.com/ryym/go-bookmark/model"

//go:generate goq

type Tables struct {
	users     model.User
	entries   model.Entry
	bookmarks model.Bookmark
}
