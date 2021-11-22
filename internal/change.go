package internal

// Change: represents the file name with full path, the number of additions and the number of removals
type Change struct {
	Added    int    `json:"added"`
	Deleted  int    `json:"removed"`
	Filename string `json:"filename"`
}
