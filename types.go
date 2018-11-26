package main

type Config struct {
	Dirs      string
	Type      string
	Format    string
	Authors   string
	From, To  string
	Output    string
	SkipEmpty bool
}

type RepoCommitCollection struct {
	Name    string    `json:"name"`
	Commits []*Commit `json:"commits"`
	MinDate int64     `json:"minDate"`
	MaxDate int64     `json:"maxDate"`
}
