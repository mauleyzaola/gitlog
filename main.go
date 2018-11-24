package main

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/golang/glog"
	"github.com/mauleyzaola/gitlog/outputs"
)

func main() {
	config := &config{
		Directory: ".",
		Type:      "commits",
		Format:    "html",
	}

	flag.StringVar(&config.Directory, "directory", config.Directory, "the path to the the git repository")
	flag.StringVar(&config.Type, "type", config.Type, "the type of output to have: [commits]")
	flag.StringVar(&config.Format, "format", config.Format, "the output format: [html|json]")
	flag.Parse()

	var (
		output   outputs.Output
		result   interface{}
		results  []interface{}
		outputFn func([]byte) error
		typeFn   func(name string, commits interface{}) (interface{}, error)
		data     []byte
		err      error
	)

	switch config.Format {
	case "json":
		output = outputs.NewJsonOutput()
	case "html":
		output = outputs.NewHTMLOutput()
	default:
		glog.Exit("unsupported format:", config.Format)
	}

	for _, repo := range strings.Fields(config.Directory) {
		switch config.Type {
		case "commits":
			typeFn = ParseCommitLines
			outputFn = output.DisplayCommits
		default:
			glog.Exit("unsupported output:", config.Type)
		}

		gitResult, err := runGitLog(repo)
		if err != nil {
			glog.Exit(err)
		}

		repoName, err := repoNameFromPath(repo)
		if err != nil {
			glog.Exit(err)
		}
		result, err = typeFn(repoName, gitResult)
		if err != nil {
			glog.Exit(err)
		}

		results = append(results, result)
	}
	if data, err = json.Marshal(&results); err != nil {
		glog.Exit(err)
	}

	if err = outputFn(data); err != nil {
		glog.Exit(err)
	}
}
