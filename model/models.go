package model

import "time"

type User struct {
	ID        int `goq:"pk"`
	Name      string
	CreatedAt time.Time
}

type Entry struct {
	ID        int `goq:"pk"`
	URL       string
	Title     string
	CreatedAt time.Time
}

type Bookmark struct {
	ID        int `goq:"pk"`
	UserID    int
	EntryID   int
	Comment   string
	CreatedAt time.Time
}
