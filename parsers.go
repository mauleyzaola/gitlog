package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// ParseLinesToCommit - Tries to convert a lines of text into a slice of Commits
// If the format is not a valid one, an error is returned
func ParseCommitLines(reader io.Reader) ([]*Commit, error) {
	scanner := bufio.NewScanner(reader)
	var (
		result []*Commit
		curr   *Commit
		ok     bool
	)
	hashes := make(map[string]struct{})
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
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
				result = append(result, curr)
			}
		}
		curr.ParseLine(fields)
	}

	return result, nil
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

// ParseLine - returns the name of the field parsed
func (t *Commit) ParseLine(fields []string) {
	if len(fields) == 0 {
		return
	}
	switch strings.ToLower(fields[0]) {
	case "commit":
		if len(fields) >= 2 {
			t.Hash = fields[1]
			return
		}
	case "author:":
		if len(fields) == 3 {
			author := &Author{
				Name:  fields[1],
				Email: fields[2],
			}
			author.Email = author.TrimEmailChars()
			t.Author = author
			return
		}
	case "authordate:":
		if len(fields) == 2 {
			value, err := time.Parse(time.RFC3339, fields[1])
			if err == nil {
				t.Date = value
				return
			}
		}
	case "commit:":
	case "commitdate:":
	default:
		ok, added, deleted := t.numStat(strings.Join(fields, " "))
		if !ok && len(t.Comment) == 0 {
			t.Comment = strings.Join(fields, " ")
			return
		} else {
			t.Added += added
			t.Deleted += deleted
			return
		}
	}
}

func (t *Commit) numStat(line string) (ok bool, added, deleted int64) {
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
