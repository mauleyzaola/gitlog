package main

import (
	"strings"
	"time"
)

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

func (t Commits) Filter(authors []string, from, to *time.Time) Commits {
	mAuth := make(map[string]struct{})
	for _, v := range authors {
		mAuth[strings.ToLower(v)] = struct{}{}
	}

	var res Commits
	for _, v := range t {
		if from != nil && v.Date.UTC().Add(time.Second).After(*from) {
			res = append(res, v)
			continue
		}
		if to != nil && v.Date.UTC().Add(-time.Second).Before(*to) {
			res = append(res, v)
			continue
		}
		if _, ok := mAuth[strings.ToLower(v.Author.Name)]; ok {
			res = append(res, v)
			continue
		}
		if _, ok := mAuth[strings.ToLower(v.Author.Email)]; ok {
			res = append(res, v)
			continue
		}
	}
	return res
}
