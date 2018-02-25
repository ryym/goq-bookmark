package repo

import "github.com/ryym/goq/goql"

func (b *Bookmarks) Entries(e *Entries) *goql.JoinDef {
	return goql.Join(e).On(e.ID.Eq(b.EntryID))
}

func (b *Bookmarks) Users(u *Users) *goql.JoinDef {
	return goql.Join(u).On(u.ID.Eq(b.UserID))
}

func (e *Entries) Bookmarks(b *Bookmarks) *goql.JoinDef {
	return goql.Join(b).On(b.EntryID.Eq(e.ID))
}
