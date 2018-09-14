package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/golang/glog"
)

func main() {
	var directory string

	flag.StringVar(&directory, "directory", "./.git", "the path to the the .git directory")
	flag.Parse()

	gitResult, err := RunGitLog(directory)
	if err != nil {
		glog.Exit(err)
	}
	commits, err := ParseCommitLines(gitResult)
	if err != nil {
		glog.Exit(err)
	}
	data, err := json.Marshal(&commits)
	if err != nil {
		glog.Exit(err)
	}
	fmt.Print(string(data))
}
