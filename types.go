package main

type config struct {
	Directory string
	Type      string
	Format    string
}

type RepoCommitCollection struct {
	Name    string    `json:"name"`
	Commits []*Commit `json:"commits"`
}
