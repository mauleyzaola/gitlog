package cmd

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/golang/glog"
	"github.com/mauleyzaola/gitlog/internal/git"
	"github.com/mauleyzaola/gitlog/internal/outputs"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generates a report output",
	Long: `report command reads from one or more git directories and outputs
HTML or JSON. For example:

gitlog report
gitlog report . --format=json
		`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlag("format", cmd.Flags().Lookup("format")); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runReportCommand(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP("format", "", "html", "path to file for storing results")
}

//nolint:gocyclo
func runReportCommand(dirs []string) error {
	config := &git.FilterParameter{
		Type:      "commits",
		Authors:   "",
		SkipEmpty: true,
	}

	if len(dirs) == 0 {
		// add current path
		dirs = []string{"."}
	}

	flag.StringVar(&config.Type, "type", config.Type, "the type of output to have: [commits]")
	flag.StringVar(&config.Authors, "authors", config.Authors, "filters by author(s)")
	flag.StringVar(&config.From, "from", config.From, "filters by start date [YYYYMMDD]")
	flag.StringVar(&config.To, "to", config.To, "filters by end date [YYYYMMDD]")
	flag.BoolVar(&config.SkipEmpty, "skip-empty", config.SkipEmpty, "skip repositories with empty data sets")

	flag.Parse()

	var (
		format     = viper.GetString("format")
		fileOutput outputs.Output
		result     interface{}
		results    []interface{}
		output     = viper.GetString("output")
		outputFn   func(*outputs.FileGenerator, interface{}) error
		typeFn     func(*git.TypeFuncParams) (bool, interface{}, error)
		err        error
		ok         bool
	)

	switch format {
	case "html":
		if viper.GetString("output") == "" {
			fileOutput = outputs.NewHTMLOutput()
		} else {
			if fileOutput, err = outputs.NewZipOutput(output); err != nil {
				return err
			}
		}
	case "json":
		fileOutput = outputs.NewJSONOutput()
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	switch config.Type {
	case "commits":
		typeFn = git.ParseCommitLines
		outputFn = fileOutput.DisplayCommits
	default:
		return fmt.Errorf("unsupported output: %s", config.Type)
	}

	started := time.Now()

	fg, err := outputs.NewFileGenerator()
	if err != nil {
		return err
	}

	repos, err := git.ParseDirNames(strings.Join(dirs, " "))
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
