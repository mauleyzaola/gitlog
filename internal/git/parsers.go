package git

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseNumber = 10
	bitSize    = 64
	commitMsg  = "commit"
)

// nolint: gocyclo
// ParseCommitLines - Tries to convert a lines of text into a slice of Commits
// If the format is not a valid one, an error is returned
// Returns true on first return parameter if there was data available
func ParseCommitLines(params *TypeFuncParams) (ok bool, result interface{}, err error) {
	// name, repoPath string, config *FilterParameter, r interface{}
	reader, ok := params.Commits.(io.Reader)
	if !ok {
		return false, nil, fmt.Errorf("cannot cast to io.Reader:%#v", params.Commits)
	}

	scanner := bufio.NewScanner(reader)
	var (
		commits          []*Commit
		curr             *Commit
		minDate, maxDate time.Time
		from, to         *time.Time
	)

	from, err = parseDate(params.Config.From)
	if err != nil {
		return false, nil, err
	}
	to, err = parseDate(params.Config.To)
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
		if hash == "" && curr == nil {
			return false, nil, fmt.Errorf("wrong git structure")
		}
		if hash != "" {
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
	if params.Config.Authors != "" {
		tmp = tmp.Filter(strings.Fields(params.Config.Authors), nil, nil)
	}

	if from != nil {
		tmp = tmp.Filter(nil, from, nil)
	}
	if to != nil {
		tmp = tmp.Filter(nil, nil, to)
	}

	files := tmp.ReadFiles(params.FullPath)

	sort.Sort(tmp)

	if len(tmp) == 0 {
		tmp = []*Commit{}
	}

	return len(tmp) != 0, &RepoCommitCollection{
		Name:     params.Name,
		Commits:  tmp,
		MinDate:  minDate.Unix(),
		MaxDate:  maxDate.Unix(),
		FileStat: tmp.FilesToMap(files),
	}, nil
}

func findHash(fields []string) string {
	const (
		hashSize      = 40
		maxFieldCount = 2
	)
	if len(fields) < maxFieldCount {
		return ""
	}
	if fields[0] != commitMsg {
		return ""
	}
	if len(fields[1]) == hashSize {
		return fields[1]
	}
	return ""
}

// ParseLine - keeps adding field data to the commit struct
func ParseLine(commit *Commit, line string) {
	const (
		hasIndex      = 1
		minFieldCount = 2
		maxFieldCount = 3
	)
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}

	log.Println("line:", line)

	switch strings.ToLower(fields[0]) {
	case commitMsg:
		if len(fields) >= minFieldCount {
			commit.Hash = fields[hasIndex]
			return
		}
	case "author:":
		if len(fields) >= maxFieldCount {
			author := &Author{
				Name:  strings.Join(fields[1:len(fields)-1], " "),
				Email: fields[len(fields)-1],
			}
			author.Email = author.TrimEmailChars()
			commit.Author = author
			return
		}
	case "authordate:":
		if len(fields) == minFieldCount {
			value, err := time.Parse(time.RFC3339, fields[1])
			if err == nil {
				commit.Date = value
				return
			}
		}
	case "commit:":
	case "commitdate:":
	default:
		ok, added, deleted, fileName := numStat(line)

		if !ok && commit.Comment == "" {
			commit.Comment = strings.Join(fields, " ")
			return
		}
		if ok {
			commit.Added += added
			commit.Deleted += deleted
			commit.Changes = append(commit.Changes, Change{Added: int(added), Deleted: int(deleted), Filename: fileName})
		}
	}
}

func numStat(line string) (ok bool, added, deleted int64, fileName string) {
	const (
		fileNameIndex = 2
		fieldCount    = 3
	)
	values := strings.Fields(line)
	if len(values) != fieldCount {
		return
	}
	var err error
	added, err = strconv.ParseInt(values[0], baseNumber, bitSize)
	if err != nil {
		return
	}
	deleted, err = strconv.ParseInt(values[1], baseNumber, bitSize)
	if err != nil {
		added = 0
		return
	}
	fileName = values[fileNameIndex]
	ok = true
	return
}

func parseDate(val string) (*time.Time, error) {
	const day = time.Hour * 24
	if val == "" {
		return nil, nil
	}
	date, err := time.Parse("20060102", val)
	if err != nil {
		return nil, err
	}
	date = date.Add(day).Add(-time.Second)
	return &date, nil
}
