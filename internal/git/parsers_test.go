package git

import (
	"bytes"
	"testing"
	"time"
)

func TestParseCommitLines(t *testing.T) {
	t.Run("single commit", commitLinesSingle)
	t.Run("multiple commits", commitLinesMultiple)
}

func commitLinesSingle(t *testing.T) {
	source := `
commit a77118ea8128202aab725841b44f919c889d949f
Author:     mauleyzaola <mauricio.leyzaola@gmail.com>
AuthorDate: 2018-08-26T01:04:55-05:00
Commit:     mauleyzaola <mauricio.leyzaola@gmail.com>
CommitDate: 2018-08-26T01:04:55-05:00

    added docker build to update chain

1	1	src/frontend/app/catalog/formula.controllers.js
`
	repoName := "unit-tests"
	buffer := bytes.NewBufferString(source)
	params := &TypeFuncParams{
		Name:    repoName,
		Commits: buffer,
	}
	_, res, err := ParseCommitLines("", nil, nil, params)
	if err != nil {
		t.Error(err)
		return
	}
	outcome, ok := res.(*RepoCommitCollection)
	if !ok {
		t.Errorf("cannot cast to *RepoCommitCollection:%#v", res)
		return
	}
	results := outcome.Commits
	if expected, actual := 1, len(results); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
		return
	}
	firstCommit := results[0]
	if expected, actual := time.Date(2018, 8, 26, 1, 4, 55, 0, time.UTC).Add(time.Hour*5).Unix(), firstCommit.Date.Unix(); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := true, firstCommit.Author != nil; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
		return
	}
	if expected, actual := "mauleyzaola", firstCommit.Author.Name; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "mauricio.leyzaola@gmail.com", firstCommit.Author.Email; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "added docker build to update chain", firstCommit.Comment; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "a77118ea8128202aab725841b44f919c889d949f", firstCommit.Hash; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commitLinesMultiple(t *testing.T) {
	source := `
commit ea8fab32b08f8d98249be02a4f0d507d75bd7dcc (HEAD -> master, origin/master, origin/HEAD)
Author:     mauleyzaola <mauricio.leyzaola@gmail.com>
AuthorDate: 2018-09-13T07:23:29-05:00
Commit:     mauleyzaola <mauricio.leyzaola@gmail.com>
CommitDate: 2018-09-13T07:23:29-05:00

    added bug report template for github

1	1	src/backend/apps/go/kaizen-consumer/main.go
8	1	src/backend/business/formula_issue.go
11	8	src/frontend/app/catalog/catalog-services.js
2	1	src/frontend/templates/issue/issue/issue.html

commit bc8d7920224b34b32578a1e95ca4159c87f19df0
Author:     Olga Acevedo <olguichi@gmail.com>
AuthorDate: 2018-09-12T20:50:35-05:00
Commit:     olguichi <olguichi@gmail.com>
CommitDate: 2018-09-12T20:50:35-05:00

    fixes #1839 - forced numeric byProduct

1	0	src/backend/business/formula.go
1	0	src/backend/interfaces/formula.go
3	0	src/frontend/templates/server/formula_print.html

commit 36f8eaeccaa1ddc23a6a09560d5319e6a87a1cf2
Author:     mauleyzaola <mauricio.leyzaola@gmail.com>
AuthorDate: 2018-08-28T18:01:18-05:00
Commit:     mauleyzaola <mauricio.leyzaola@gmail.com>
CommitDate: 2018-08-28T18:01:18-05:00

    fixes #1828 - automate service restart

30	0	.github/ISSUE_TEMPLATE/bug_report.md

`
	buffer := bytes.NewBufferString(source)
	repoName := "unit-tests"
	params := &TypeFuncParams{
		Name:    repoName,
		Commits: buffer,
	}
	_, res, err := ParseCommitLines("", nil, nil, params)
	if err != nil {
		t.Error(err)
		return
	}
	outcome, ok := res.(*RepoCommitCollection)
	if !ok {
		t.Errorf("cannot cast to *RepoCommitCollection:%#v", res)
		return
	}
	results := outcome.Commits
	samples := []Commit{
		{
			Date:    time.Date(2018, 8, 28, 18, 1, 18, 0, time.UTC).Add(time.Hour * 5),
			Author:  &Author{Name: "mauleyzaola", Email: "mauricio.leyzaola@gmail.com"},
			Hash:    "36f8eaeccaa1ddc23a6a09560d5319e6a87a1cf2",
			Comment: "fixes #1828 - automate service restart",
			Added:   30,
			Deleted: 0,
		},
		{
			Date:    time.Date(2018, 9, 12, 20, 50, 35, 0, time.UTC).Add(time.Hour * 5),
			Author:  &Author{Name: "Olga Acevedo", Email: "olguichi@gmail.com"},
			Hash:    "bc8d7920224b34b32578a1e95ca4159c87f19df0",
			Comment: "fixes #1839 - forced numeric byProduct",
			Added:   5,
			Deleted: 0,
		},
		{
			Date:    time.Date(2018, 9, 13, 7, 23, 29, 0, time.UTC).Add(time.Hour * 5),
			Author:  &Author{Name: "mauleyzaola", Email: "mauricio.leyzaola@gmail.com"},
			Hash:    "ea8fab32b08f8d98249be02a4f0d507d75bd7dcc",
			Comment: "added bug report template for github",
			Added:   22,
			Deleted: 11,
		},
	}

	if expected, actual := len(samples), len(results); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
		return
	}

	for i := 0; i < len(samples); i++ {
		sample := samples[i]
		result := results[i]
		if expected, actual := sample.Date.Unix(), result.Date.Unix(); expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := sample.Author != nil, result.Author != nil; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
			continue
		}
		if expected, actual := sample.Author.Name, result.Author.Name; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := sample.Author.Email, result.Author.Email; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := sample.Comment, result.Comment; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := sample.Hash, result.Hash; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := sample.Added, result.Added; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := sample.Deleted, result.Deleted; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
	}
}

func TestCommit_ParseLine(t *testing.T) {
	t.Run("author", commitParseLineAuthor)
	t.Run("author date", commitAuthorDate)
	t.Run("commit hash", commitHash)
	t.Run("blank", commitBlank)
	t.Run("comment", commitComment)
	t.Run("changes", commitChange)
}

func commitParseLineAuthor(t *testing.T) {
	c := &Commit{}
	line := "Author:     mauleyzaola <mauricio.leyzaola@gmail.com>"
	ParseLine(c, line)
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

func commitChange(t *testing.T) {
	c := &Commit{}
	line := "235     4       commonpasswords/main.go"
	ParseLine(c, line)

	if expected, actual := 1, len(c.Changes); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 235, c.Changes[0].Added; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := 4, c.Changes[0].Deleted; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "commonpasswords/main.go", c.Changes[0].Filename; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commitHash(t *testing.T) {
	c := &Commit{}
	line := "commit a77118ea8128202aab725841b44f919c889d949f"
	ParseLine(c, line)
	if expected, actual := "a77118ea8128202aab725841b44f919c889d949f", c.Hash; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commitAuthorDate(t *testing.T) {
	c := &Commit{}
	line := "AuthorDate: 2018-08-26T01:04:55-05:00"
	ParseLine(c, line)
	if expected, actual := time.Date(2018, 8, 26, 1, 4, 55, 0, time.UTC).Add(time.Hour*5).Unix(), c.Date.Unix(); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commitBlank(t *testing.T) {
	c := &Commit{}
	ParseLine(c, "")
	if expected, actual := "", c.Hash; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := true, c.Date.IsZero(); expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := true, c.Author == nil; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := "", c.Comment; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := int64(0), c.Added; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
	if expected, actual := int64(0), c.Deleted; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func commitComment(t *testing.T) {
	c := &Commit{}
	line := "added docker build to update chain"
	ParseLine(c, line)
	if expected, actual := line, c.Comment; expected != actual {
		t.Errorf("expected:%v actual:%v", expected, actual)
	}
}

func Test_IsNumStat(t *testing.T) {
	cases := []struct {
		expected       bool
		added, deleted int64
		line           string
	}{
		{
			expected: true,
			added:    1,
			deleted:  1,
			line: "1	1	src/frontend/app/catalog/formula.controllers.js",
		},
		{
			expected: false,
			added:    0,
			deleted:  0,
			line:     "fixes #1828 - automate service restart",
		},
		{
			expected: true,
			added:    11,
			deleted:  8,
			line: "11	8	src/frontend/app/catalog/catalog-services.js",
		},
	}

	for i, v := range cases {
		ok, added, deleted, _ := numStat(v.line)
		if expected, actual := v.expected, ok; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := v.added, added; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
		if expected, actual := v.deleted, deleted; expected != actual {
			t.Errorf("[%d] - expected:%v actual:%v", i, expected, actual)
		}
	}
}
