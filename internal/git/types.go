package git

type RepoCommitCollection struct {
	Name    string    `json:"name"`
	Commits []*Commit `json:"commits"`
	MinDate int64     `json:"minDate"`
	MaxDate int64     `json:"maxDate"`
	// each key is the extension and the value is the sum of its sizes and the count of each file type
	FileStat map[string]*RepoFileInfo `json:"fileStat"`
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
	Name     string
	FullPath string
	Commits  interface{}
}
