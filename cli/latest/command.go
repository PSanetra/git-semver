package latest

import (
	"fmt"
	"github.com/psanetra/git-semver/cli/common_opts"
	"github.com/psanetra/git-semver/latest"
	"github.com/psanetra/git-semver/logger"
	"github.com/spf13/cobra"
)

var includePreReleases bool

var Command = cobra.Command{
	Use:   "latest",
	Short: "prints latest semantic version",
	Long:  `This command prints the latest semantic version in the current repository by comparing all git tags. Tag names may have a "v" prefix, but this commands prints the version always without that prefix.`,
	Run: func(cmd *cobra.Command, args []string) {

		latestVersion, err := latest.Latest(latest.LatestOptions{
			Workdir:            common_opts.Workdir,
			IncludePreReleases: includePreReleases,
		})

		if err != nil {
			logger.Logger.Fatalln(err)
		}

		fmt.Print(latestVersion.ToString())

	},
}

func init() {
	Command.Flags().BoolVar(&includePreReleases, "include-pre-releases", false, "Also consider pre-releases as the latest version")
}
