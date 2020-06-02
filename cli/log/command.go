package log

import (
	"fmt"
	"github.com/psanetra/git-semver/cli/common_opts"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
	"github.com/psanetra/git-semver/version_log"
	"github.com/spf13/cobra"
)

var excludePreReleases bool

var Command = cobra.Command{
	Use:   "log [<version>]",
	Short: "prints the git log for the specified version",
	Long: "This command prints all commits, which were contained in a specified version or all commits since the latest version if no version is specified.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var version *semver.Version
		var err error

		if len(args) > 0 {
			version, err = semver.ParseVersion(args[0])
			if err != nil {
				logger.Logger.Fatalln("Could not parse version:", err)
			}
		}

		commits, err := version_log.VersionLog(version_log.VersionLogOptions{
			Workdir:                  common_opts.Workdir,
			Version:                  version,
			ExcludePreReleaseCommits: excludePreReleases,
		})

		if err != nil {
			logger.Logger.Fatalln(err)
		}

		for _, commit := range commits{
			fmt.Print(commit)
		}
	},
}

func init() {
	Command.Flags().BoolVar(&excludePreReleases, "exclude-pre-releases", false, "Specifies if the log should exclude pre-release commits from the log.")
}
