package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/mauleyzaola/gitlog/outputs"
)

func main() {
	config := &Config{
		Dirs:    ".",
		Type:    "commits",
		Format:  "html",
		Authors: "",
		Output:  "",
	}

	flag.StringVar(&config.Dirs, "dirs", config.Dirs, "the path(s) to the the git repository")
	flag.StringVar(&config.Type, "type", config.Type, "the type of output to have: [commits]")
	flag.StringVar(&config.Format, "format", config.Format, "the output format: [html|json]")
	flag.StringVar(&config.Authors, "authors", config.Authors, "filters by author(s)")
	flag.StringVar(&config.From, "from", config.From, "filters by start date [YYYYMMDD]")
	flag.StringVar(&config.To, "to", config.To, "filters by end date [YYYYMMDD]")
	flag.StringVar(&config.Output, "output", config.Output, "path to zip file for storing results")

	flag.Parse()

	var (
		output   outputs.Output
		result   interface{}
		results  []interface{}
		outputFn func(*outputs.FileGenerator, interface{}) error
		typeFn   func(name string, params *Config, commits interface{}) (interface{}, error)
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

	switch config.Type {
	case "commits":
		typeFn = ParseCommitLines
		outputFn = output.DisplayCommits
	default:
		glog.Exit("unsupported output:", config.Type)
	}

	fg, err := outputs.NewFileGenerator()
	if err != nil {
		glog.Exitln(err)
	}

	repos, err := parseDirNames(config.Dirs)
	if err != nil {
		glog.Exitln(err)
	}
	for _, repo := range repos {
		gitResult, err := runGitLog(repo)
		if err != nil {
			glog.Warningf("cannot obtain git log information from directory:%s. %s", repo, err)
			continue
		}

		repoName, err := repoNameFromPath(repo)
		if err != nil {
			glog.Exit(err)
		}
		result, err = typeFn(repoName, config, gitResult)
		if err != nil {
			glog.Exit(err)
		}

		results = append(results, result)
	}

	if err = outputFn(fg, results); err != nil {
		glog.Exit(err)
	}
}
