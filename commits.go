package main

import (
	"os"
	"path/filepath"
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
		if _, ok := mAuth[strings.ToLower(v.Author.Email)]; ok {
			res = append(res, v)
			continue
		}
	}
	return res
}

func (t Commits) ReadFiles(repoPath string) []*RepoFile {
	var (
		res  []*RepoFile
		keys = make(map[string]struct{})
	)
	for _, commit := range t {
		for _, change := range commit.Changes {
			file := filepath.Join(repoPath, change.Filename)
			// ignore hidden files
			if strings.HasPrefix(filepath.Base(file), ".") {
				continue
			}
			if _, ok := keys[file]; ok {
				// avoid duplicates
				continue
			}
			keys[file] = struct{}{}
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			res = append(res, &RepoFile{
				Name: file,
				Size: info.Size(),
			})
		}
	}
	return res
}

func (t Commits) FilesToMap(files []*RepoFile) map[string]*RepoFileInfo {
	res := make(map[string]*RepoFileInfo)
	for _, f := range files {
		ext := filepath.Ext(f.Name)
		if len(ext) == 0 {
			continue
		}
		val, ok := res[ext]
		if !ok {
			val = &RepoFileInfo{}
		}
		val.Count++
		val.Size += f.Size
		res[ext] = val
	}
	return res
}
