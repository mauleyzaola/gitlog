package main

import "time"

// Commit - Type that holds the information about one single commit
type Commit struct {
	Hash    string    `json:"hash"`
	Author  *Author   `json:"author"`
	Date    time.Time `json:"date"` // for the sake of simplicity we deal only with AuthorDate
	Comment string    `json:"-"`    // issues with JSON.parse(), need more rework
	Added   int64     `json:"added"`
	Deleted int64     `json:"removed"`
	Changes []Change  `json:"changes"`
}
