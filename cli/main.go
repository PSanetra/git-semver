package main

import (
	"github.com/psanetra/git-semver/cli/common_opts"
	"github.com/psanetra/git-semver/cli/compare"
	"github.com/psanetra/git-semver/cli/latest"
	"github.com/psanetra/git-semver/cli/next"
	"github.com/psanetra/git-semver/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-semver",
	Short: "git-semver is a cli tool to apply semver conventions to git based projects.",
	// Long: `git-semver is a cli tool to apply semver conventions to git based projects.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		processLogLevelFlag(cmd)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Fatalln(err)
	}
}

func main() {

	rootCmd.AddCommand(&latest.Command)
	rootCmd.AddCommand(&next.Command)
	rootCmd.AddCommand(&compare.Command)
	err := rootCmd.Execute()

	if err != nil {
		logger.Logger.Fatalln(err)
	}

}

func init() {
	rootCmd.PersistentFlags().StringVarP(&common_opts.Workdir, "workdir", "w", ".", "Working directory to use")
	rootCmd.PersistentFlags().String("log-level", logger.DEFAULT_LOG_LEVEL.String(), "panic | fatal | error | warn | info | debug | trace")
}

func processLogLevelFlag(cmd *cobra.Command) {
	logLevel := cmd.Flag("log-level").Value.String()
	logger.SetLevel(logLevel)
}
