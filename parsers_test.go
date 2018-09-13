package gitlog

import (
	"bytes"
	"testing"
	"time"
)

func TestParseLinesToCommit(t *testing.T) {
	t.Skip()
	source := `
commit a77118ea8128202aab725841b44f919c889d949f
Author:     mauleyzaola <mauricio.leyzaola@gmail.com>
AuthorDate: 2018-08-26T01:04:55-05:00
Commit:     mauleyzaola <mauricio.leyzaola@gmail.com>
CommitDate: 2018-08-26T01:04:55-05:00

    added docker build to update chain
`
	buffer := bytes.NewBufferString(source)
	result, err := ParseLinesToCommits(buffer)
	if err != nil {
		t.Error(err)
		return
	}
	if expected, actual := 1, len(result); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func TestCommit_ParseLine(t *testing.T) {
	t.Run("author", commitParseLineAuthor)
	t.Run("author date", commit_parseLine5)
	t.Run("commit hash", commit_parseLineHash)
	t.Run("blank", commit_parseLineBlank)
	t.Run("comment", commit_parseLineComment)
}

func commitParseLineAuthor(t *testing.T) {
	c := &Commit{}
	line := "Author:     mauleyzaola <mauricio.leyzaola@gmail.com>"
	ok := c.ParseLine(line)
	if expected, actual := true, ok; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if c.Author == nil {
		t.Error("author is nil")
		return
	}
	if expected, actual := "mauleyzaola", c.Author.Name; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "mauricio.leyzaola@gmail.com", c.Author.Email; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commit_parseLineHash(t *testing.T) {
	c := &Commit{}
	line := "commit a77118ea8128202aab725841b44f919c889d949f"
	ok := c.ParseLine(line)
	if expected, actual := true, ok; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "a77118ea8128202aab725841b44f919c889d949f", c.Hash; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commit_parseLine5(t *testing.T) {
	c := &Commit{}
	line := "AuthorDate: 2018-08-26T01:04:55-05:00"
	ok := c.ParseLine(line)
	if expected, actual := true, ok; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := time.Date(2018, 8, 26, 1, 4, 55, 0, time.UTC).Add(time.Hour*5).Unix(), c.Date.Unix(); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commit_parseLineBlank(t *testing.T) {
	c := &Commit{}
	ok := c.ParseLine("")
	if expected, actual := false, ok; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commit_parseLineComment(t *testing.T) {
	c := &Commit{}
	line := "added docker build to update chain"
	ok := c.ParseLine(line)
	if expected, actual := true, ok; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := line, c.Comment; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}
