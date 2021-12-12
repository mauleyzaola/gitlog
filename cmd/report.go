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
		var flagNames = []string{
			"authors",
			"format",
			"from",
			"to",
			"type",
		}
		for _, v := range flagNames {
			if err := viper.BindPFlag(v, cmd.Flags().Lookup(v)); err != nil {
				return err
			}
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
	reportCmd.Flags().StringP("authors", "", "", "filter commits by the author(s)")
	reportCmd.Flags().StringP("format", "", "html", "path to file for storing results")
	reportCmd.Flags().StringP("from", "", "", "filters by start date [YYYYMMDD]")
	reportCmd.Flags().StringP("to", "", "", "filters by end date [YYYYMMDD]")
	reportCmd.Flags().StringP("type", "", "commits", "type of output [commits]")
}

//nolint:gocyclo
func runReportCommand(dirs []string) error {
	config := &git.FilterParameter{
		SkipEmpty: true,
	}

	if len(dirs) == 0 {
		// add current path
		dirs = []string{"."}
	}

	flag.BoolVar(&config.SkipEmpty, "skip-empty", config.SkipEmpty, "skip repositories with empty data sets")

	flag.Parse()

	var (
		authors    = viper.GetString("authors")
		format     = viper.GetString("format")
		fileOutput outputs.Output
		from, to   *time.Time
		result     interface{}
		results    []interface{}
		output     = viper.GetString("output")
		outputFn   func(*outputs.FileGenerator, interface{}) error
		typeName   = viper.GetString("type")
		typeFn     func(
			authors string,
			from, to *time.Time,
			params *git.TypeFuncParams,
		) (bool, interface{}, error)
		err error
		ok  bool
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

	switch typeName {
	case "commits":
		typeFn = git.ParseCommitLines
		outputFn = fileOutput.DisplayCommits
	default:
		return fmt.Errorf("unsupported output: %s", typeName)
	}

	if val := viper.GetString("from"); val != "" {
		if from, err = parseDate(val); err != nil {
			return err
		}
	}
	if val := viper.GetString("to"); val != "" {
		if to, err = parseDate(val); err != nil {
			return err
		}
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
		ok, result, errGl = typeFn(
			authors,
			from, to,
			&git.TypeFuncParams{
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

// TODO: move this functions somewhere else
func parseDate(val string) (*time.Time, error) {
	const day = time.Hour * 24
	if val == "" {
		return nil, nil
	}
	date, err := time.Parse("20060102", val)
	if err != nil {
		return nil, err
	}
	date = date.Add(day).Add(-time.Second)
	return &date, nil
}
