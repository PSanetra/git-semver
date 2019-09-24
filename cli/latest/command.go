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
	Short: "Prints latest semantic version",
	Run: func(cmd *cobra.Command, args []string) {

		nextVersion, err := latest.Latest(latest.LatestOptions{
			Workdir: common_opts.Workdir,
			IncludePreReleases: includePreReleases,
		})

		if err != nil {
			logger.Logger.Fatalln(err)
		}

		fmt.Print(nextVersion.ToString())

	},
}

func init() {
	Command.Flags().BoolVar(&includePreReleases, "include-pre-releases", false, "Also consider pre-releases as the latest version")
}
