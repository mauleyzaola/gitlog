package cmd

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang/glog"
	"github.com/mauleyzaola/gitlog/internal/git"
	"github.com/mauleyzaola/gitlog/internal/outputs"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runReportCommand(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runReportCommand() error {
	config := &git.FilterParameter{
		Dirs:      ".",
		Type:      "commits",
		Format:    "html",
		Authors:   "",
		Output:    "",
		SkipEmpty: true,
	}

	flag.StringVar(&config.Dirs, "dirs", config.Dirs, "the path(s) to the the git repository")
	flag.StringVar(&config.Type, "type", config.Type, "the type of output to have: [commits]")
	flag.StringVar(&config.Format, "format", config.Format, "the output format: [html|json]")
	flag.StringVar(&config.Authors, "authors", config.Authors, "filters by author(s)")
	flag.StringVar(&config.From, "from", config.From, "filters by start date [YYYYMMDD]")
	flag.StringVar(&config.To, "to", config.To, "filters by end date [YYYYMMDD]")
	flag.StringVar(&config.Output, "output", config.Output, "path to file for storing results")
	flag.BoolVar(&config.SkipEmpty, "skip-empty", config.SkipEmpty, "skip repositories with empty data sets")

	flag.Parse()

	var (
		output   outputs.Output
		result   interface{}
		results  []interface{}
		outputFn func(*outputs.FileGenerator, interface{}) error
		typeFn   func(*git.TypeFuncParams) (bool, interface{}, error)
		err      error
		ok       bool
	)

	if config.Format == "json" {
		output = outputs.NewJSONOutput()
	} else if config.Format == "html" {
		if config.Output == "" {
			output = outputs.NewHTMLOutput()
		} else {
			if output, err = outputs.NewZipOutput(config.Output); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unsupported format: %s", config.Format)
	}

	switch config.Type {
	case "commits":
		typeFn = git.ParseCommitLines
		outputFn = output.DisplayCommits
	default:
		return fmt.Errorf("unsupported output: %s", config.Type)
	}

	started := time.Now()

	fg, err := outputs.NewFileGenerator()
	if err != nil {
		return err
	}

	repos, err := git.ParseDirNames(config.Dirs)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		gitResult, errGl := git.RunGitLog(repo)
		if errGl != nil {
			glog.Warningf("cannot obtain git log information from directory:%s. %s", repo, errGl)
			continue
		}

		repoName, errGl := git.RepoNameFromPath(repo)
		if errGl != nil {
			return errGl
		}
		ok, result, errGl = typeFn(&git.TypeFuncParams{
			Config:   config,
			Name:     repoName,
			FullPath: repo,
			Commits:  gitResult,
		})
		if errGl != nil {
			glog.Exit(errGl)
		}

		if ok || !config.SkipEmpty {
			results = append(results, result)
		}
	}

	if errOutput := outputFn(fg, results); errOutput != nil {
		return errOutput
	}

	log.Println("total time elapsed:", time.Since(started))
	return nil
}
