package repo

import "github.com/ryym/goq/gql"

func (b *Bookmarks) Entries(e *Entries) *gql.JoinDef {
	return gql.Join(e).On(b.EntryID.Eq(e.ID))
}
