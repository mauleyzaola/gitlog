package main

type Commits []*Commit

func (t Commits) Len() int {
	return len(t)
}

func (t Commits) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Commits) Less(i, j int) bool {
	return t[i].Date.Before(t[j].Date)
}
