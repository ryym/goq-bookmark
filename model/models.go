package model

type User struct {
	ID   int `goq:"pk"`
	Name string
}

type Entry struct {
	ID    int `goq:"pk"`
	URL   string
	Title string
}

type Bookmark struct {
	ID      int `goq:"pk"`
	UserID  int
	EntryID int
	Comment string
}
