// DO NOT EDIT. This code is generated by Goq.
// https://github.com/ryym/goq

package repo

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/dialect"
)

type Users struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID        *goq.Column
	Name      *goq.Column
	CreatedAt *goq.Column
}

func NewUsers(alias string) *Users {
	cm := goq.NewColumnMaker("User", "users").As(alias)
	t := &Users{

		ID:        cm.Col("ID", "id").PK().Bld(),
		Name:      cm.Col("Name", "name").Bld(),
		CreatedAt: cm.Col("CreatedAt", "created_at").Bld(),
	}
	cols := []*goq.Column{t.ID, t.Name, t.CreatedAt}
	t.Table = goq.NewTable("users", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Users) As(alias string) *Users { return NewUsers(alias) }

type Entries struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID        *goq.Column
	URL       *goq.Column
	Title     *goq.Column
	CreatedAt *goq.Column
}

func NewEntries(alias string) *Entries {
	cm := goq.NewColumnMaker("Entry", "entries").As(alias)
	t := &Entries{

		ID:        cm.Col("ID", "id").PK().Bld(),
		URL:       cm.Col("URL", "url").Bld(),
		Title:     cm.Col("Title", "title").Bld(),
		CreatedAt: cm.Col("CreatedAt", "created_at").Bld(),
	}
	cols := []*goq.Column{t.ID, t.URL, t.Title, t.CreatedAt}
	t.Table = goq.NewTable("entries", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Entries) As(alias string) *Entries { return NewEntries(alias) }

type Bookmarks struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID        *goq.Column
	UserID    *goq.Column
	EntryID   *goq.Column
	Comment   *goq.Column
	CreatedAt *goq.Column
}

func NewBookmarks(alias string) *Bookmarks {
	cm := goq.NewColumnMaker("Bookmark", "bookmarks").As(alias)
	t := &Bookmarks{

		ID:        cm.Col("ID", "id").PK().Bld(),
		UserID:    cm.Col("UserID", "user_id").Bld(),
		EntryID:   cm.Col("EntryID", "entry_id").Bld(),
		Comment:   cm.Col("Comment", "comment").Bld(),
		CreatedAt: cm.Col("CreatedAt", "created_at").Bld(),
	}
	cols := []*goq.Column{t.ID, t.UserID, t.EntryID, t.Comment, t.CreatedAt}
	t.Table = goq.NewTable("bookmarks", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Bookmarks) As(alias string) *Bookmarks { return NewBookmarks(alias) }

type Builder struct {
	*goq.Builder

	Users     *Users
	Entries   *Entries
	Bookmarks *Bookmarks
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),

		Users:     NewUsers(""),
		Entries:   NewEntries(""),
		Bookmarks: NewBookmarks(""),
	}
}
