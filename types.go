package main

type Config struct {
	Directories string
	Type        string
	Format      string
	Authors     string
	From, To    string
}

type RepoCommitCollection struct {
	Name    string    `json:"name"`
	Commits []*Commit `json:"commits"`
	MinDate int64     `json:"minDate"`
	MaxDate int64     `json:"maxDate"`
}
