package gitlog

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// ParseLinesToCommit - Tries to convert a lines of text into a slice of Commits
// If the format is not a valid one, an error is returned
func ParseCommitLines(reader io.Reader) ([]Commit, error) {
	scanner := bufio.NewScanner(reader)
	var (
		result []Commit
		curr   *Commit
	)
	for scanner.Scan() {
		line := scanner.Text()
		if curr == nil {
			curr = &Commit{}
		}
		if ok := curr.ParseLine(line); !ok {
			// TODO: deal with commits without comments
			if len(curr.Comment) != 0 {
				result = append(result, *curr)
				curr = nil
			}
		}
	}
	if curr != nil {
		result = append(result, *curr)
	}
	return result, nil
}

// ParseLine - returns true if the line value affects any field
func (t *Commit) ParseLine(line string) bool {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false
	}
	switch strings.ToLower(fields[0]) {
	case "commit":
		if len(fields) >= 2 {
			t.Hash = fields[1]
			return true
		}
	case "author:":
		if len(fields) == 3 {
			author := &Author{
				Name:  fields[1],
				Email: fields[2],
			}
			author.Email = author.TrimEmailChars()
			t.Author = author
			return true
		}
	case "authordate:":
		if len(fields) == 2 {
			value, err := time.Parse(time.RFC3339, fields[1])
			if err != nil {
				return false
			}
			t.Date = value
			return true
		}
	case "commit:":
		fallthrough
	case "commitdate:":
		return false
	default:
		t.Comment = strings.Join(fields, " ")
		return true
	}
	return false
}
