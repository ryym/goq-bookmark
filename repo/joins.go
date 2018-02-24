package repo

import "github.com/ryym/goq/gql"

func (b *Bookmarks) Entries(e *Entries) *gql.JoinDef {
	return gql.Join(e).On(e.ID.Eq(b.EntryID))
}

func (b *Bookmarks) Users(u *Users) *gql.JoinDef {
	return gql.Join(u).On(u.ID.Eq(b.UserID))
}

func (e *Entries) Bookmarks(b *Bookmarks) *gql.JoinDef {
	return gql.Join(b).On(b.EntryID.Eq(e.ID))
}
