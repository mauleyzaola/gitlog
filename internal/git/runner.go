package git

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunGitLog(directory string) (io.Reader, error) {
	params := []string{fmt.Sprintf("--git-dir=%s", filepath.Join(directory, ".git"))}
	params = append(params, strings.Fields("log --no-merges  --pretty=fuller --date=iso-strict --numstat")...)

	cmd := exec.Command("git", params...)
	stdErr := &bytes.Buffer{}
	stdOut := &bytes.Buffer{}
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut
	if err := cmd.Run(); err != nil {
		log.Println(stdErr.String())
		return nil, err
	}
	return stdOut, nil
}

func RepoNameFromPath(p string) (string, error) {
	if p == "." {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Base(dir), nil
	}
	return filepath.Base(p), nil
}

func ParseDirNames(dirs string) ([]string, error) {
	var res []string
	for _, dir := range strings.Fields(dirs) {
		values, err := filepath.Glob(dir)
		if err != nil {
			return nil, err
		}
		res = append(res, values...)
	}
	return res, nil
}
