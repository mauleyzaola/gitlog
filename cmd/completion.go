package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var supportedShells = []string{
	"bash",
	"fish",
	"power-shell",
	"zsh",
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use: "completion",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
	Short: "Generates completion scripts for a shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := startCompletion(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func startCompletion(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing shell, please choose any of these: %s", supportedShells)
	}
	shell := strings.ToLower(args[0])
	var output = os.Stdout
	switch shell {
	case "bash":
		return rootCmd.GenBashCompletion(output)
	case "fish":
		return rootCmd.GenFishCompletion(output, true)
	case "power-shell":
		return rootCmd.GenPowerShellCompletion(output)
	case "zsh":
		return rootCmd.GenZshCompletion(output)
	default:
		return fmt.Errorf("unsupported shell: %s, please choose any of these: %s", shell, supportedShells)
	}
}
