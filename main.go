package main

import (
	"encoding/json"
	"flag"

	"github.com/golang/glog"
	"github.com/mauleyzaola/gitlog/outputs"
)

func main() {
	config := &config{
		Directory: "./.git",
		Type:      "commits",
		Format:    "html",
	}

	flag.StringVar(&config.Directory, "directory", config.Directory, "the path to the the .git directory")
	flag.StringVar(&config.Type, "type", config.Type, "the type of output to have: [commits]")
	flag.StringVar(&config.Format, "format", config.Format, "the output format: [html|json]")
	flag.Parse()

	var (
		output   outputs.Output
		result   interface{}
		outputFn func([]byte) error
		typeFn   func(interface{}) (interface{}, error)
	)

	switch config.Format {
	case "json":
		output = outputs.NewJsonOutput()
	case "html":
		output = outputs.NewHTMLOutput()
	default:
		glog.Exit("unsupported format:", config.Format)
	}

	switch config.Type {
	case "commits":
		typeFn = ParseCommitLines
		outputFn = output.DisplayCommits
	default:
		glog.Exit("unsupported output:", config.Type)
	}

	gitResult, err := runGitLog(config.Directory)
	if err != nil {
		glog.Exit(err)
	}

	result, err = typeFn(gitResult)
	if err != nil {
		glog.Exit(err)
	}

	data, err := json.Marshal(&result)
	if err != nil {
		glog.Exit(err)
	}

	if err = outputFn(data); err != nil {
		glog.Exit(err)
	}
}
