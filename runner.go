package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

func runGitLog(directory string) (io.Reader, error) {
	params := []string{fmt.Sprintf("--git-dir=%s", filepath.Join(directory, ".git"))}
	params = append(params, strings.Fields("log --no-merges  --pretty=fuller --date=iso-strict --numstat")...)

	cmd := exec.Command("git", params...)
	stdErr := &bytes.Buffer{}
	stdOut := &bytes.Buffer{}
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut
	if err := cmd.Run(); err != nil {
		glog.Error(stdErr.String())
		return nil, err
	}
	return stdOut, nil
}
