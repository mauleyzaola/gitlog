package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ParseCommitLines - Tries to convert a lines of text into a slice of Commits
// If the format is not a valid one, an error is returned
func ParseCommitLines(name string, r interface{}) (interface{}, error) {
	reader, ok := r.(io.Reader)
	if !ok {
		return nil, fmt.Errorf("cannot cast to io.Reader:%#v", r)
	}
	scanner := bufio.NewScanner(reader)
	var (
		commits          []*Commit
		curr             *Commit
		minDate, maxDate time.Time
	)
	hashes := make(map[string]struct{})
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			// blank line, ignore
			continue
		}
		hash := findHash(fields)
		if len(hash) == 0 && curr == nil {
			return nil, fmt.Errorf("wrong git structure")
		}
		if len(hash) != 0 {
			if _, ok = hashes[hash]; !ok {
				// new commit detected
				curr = &Commit{Hash: hash}
				hashes[hash] = struct{}{}
				commits = append(commits, curr)
			}
		}
		ParseLine(curr, line)
		if minDate.IsZero() || curr.Date.Before(minDate) {
			minDate = curr.Date
		}
		if maxDate.Before(curr.Date) {
			maxDate = curr.Date
		}
	}
	tmp := Commits(commits)
	sort.Sort(tmp)

	return &RepoCommitCollection{
		Name:    name,
		Commits: commits,
		MinDate: minDate.Unix(),
		MaxDate: maxDate.Unix(),
	}, nil
}

func findHash(fields []string) string {
	if len(fields) < 2 {
		return ""
	}
	if fields[0] != "commit" {
		return ""
	}
	if len(fields[1]) == 40 {
		return fields[1]
	}
	return ""
}

// ParseLine - keeps adding field data to the commit struct
func ParseLine(commit *Commit, line string) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}
	switch strings.ToLower(fields[0]) {
	case "commit":
		if len(fields) >= 2 {
			commit.Hash = fields[1]
			return
		}
	case "author:":
		if len(fields) >= 3 {
			author := &Author{
				Name:  strings.Join(fields[1:len(fields)-1], " "),
				Email: fields[len(fields)-1],
			}
			author.Email = author.TrimEmailChars()
			commit.Author = author
			return
		}
	case "authordate:":
		if len(fields) == 2 {
			value, err := time.Parse(time.RFC3339, fields[1])
			if err == nil {
				commit.Date = value
				return
			}
		}
	case "commit:":
	case "commitdate:":
	default:
		ok, added, deleted := numStat(line)
		if !ok && len(commit.Comment) == 0 {
			commit.Comment = strings.Join(fields, " ")
			return
		} else if ok {
			commit.Added += added
			commit.Deleted += deleted
			return
		}
	}
}

func numStat(line string) (ok bool, added, deleted int64) {
	values := strings.Split(line, "\t")
	if len(values) != 3 {
		return
	}
	var err error
	added, err = strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return
	}
	deleted, err = strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		added = 0
		return
	}
	ok = true
	return
}
