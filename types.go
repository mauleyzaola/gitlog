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
	Name     string                   `json:"name"`
	Commits  []*Commit                `json:"commits"`
	MinDate  int64                    `json:"minDate"`
	MaxDate  int64                    `json:"maxDate"`
	FileStat map[string]*RepoFileInfo `json:"fileStat"` // each key is the extension and the value is the sum of its sizes and the count of each file type
}

type RepoFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type RepoFileInfo struct {
	Size  int64 `json:"size"`
	Count int   `json:"count"`
}

type TypeFuncParams struct {
	config   *Config
	name     string
	fullPath string
	params   *Config
	commits  interface{}
}
