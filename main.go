package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/golang/glog"
)

func main() {
	config := &config{
		Directory: "./.git",
		Output:    "commits",
	}

	flag.StringVar(&config.Directory, "directory", config.Directory, "the path to the the .git directory")
	flag.StringVar(&config.Output, "output", config.Output, "the type of output to have: [commits]")
	flag.Parse()

	switch config.Output {
	case "commits":
	default:
		glog.Exit("unsupported output:", config.Output)
	}

	gitResult, err := RunGitLog(config.Directory)
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
