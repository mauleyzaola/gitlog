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
// Returns true on first return parameter if there was data available
func ParseCommitLines(name string, config *Config, r interface{}) (bool, interface{}, error) {
	reader, ok := r.(io.Reader)
	if !ok {
		return false, nil, fmt.Errorf("cannot cast to io.Reader:%#v", r)
	}

	scanner := bufio.NewScanner(reader)
	var (
		commits          []*Commit
		curr             *Commit
		minDate, maxDate time.Time
		from, to         *time.Time
	)

	from, err := parseDate(config.From)
	if err != nil {
		return false, nil, err
	}
	to, err = parseDate(config.To)
	if err != nil {
		return false, nil, err
	}

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
			return false, nil, fmt.Errorf("wrong git structure")
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

	// cast to type
	tmp := Commits(commits)

	// apply filters
	if len(config.Authors) != 0 {
		tmp = tmp.Filter(strings.Fields(config.Authors), nil, nil)
	}

	if from != nil {
		tmp = tmp.Filter(nil, from, nil)
	}
	if to != nil {
		tmp = tmp.Filter(nil, nil, to)
	}

	sort.Sort(tmp)

	if len(tmp) == 0 {
		tmp = []*Commit{}
	}

	return len(tmp) != 0, &RepoCommitCollection{
		Name:    name,
		Commits: tmp,
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

func parseDate(val string) (*time.Time, error) {
	if len(val) == 0 {
		return nil, nil
	}
	date, err := time.Parse("20060102", val)
	if err != nil {
		return nil, err
	}
	date = date.Add(time.Hour * 24).Add(-time.Second)
	return &date, nil
}
