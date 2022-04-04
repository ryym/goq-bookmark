package repo

import "github.com/ryym/goq"

func (b *Bookmarks) Entries(e *Entries) *goq.JoinDef {
	return goq.Join(e).On(e.ID.Eq(b.EntryID))
}

func (b *Bookmarks) Users(u *Users) *goq.JoinDef {
	return goq.Join(u).On(u.ID.Eq(b.UserID))
}

func (e *Entries) Bookmarks(b *Bookmarks) *goq.JoinDef {
	return goq.Join(b).On(b.EntryID.Eq(e.ID))
}
