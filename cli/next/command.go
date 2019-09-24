package next

import (
	"fmt"
	"github.com/psanetra/git-semver/cli/common_opts"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/next"
	"github.com/psanetra/git-semver/semver"
	"github.com/spf13/cobra"
)

var projectName string
var stable bool
var preReleaseTag string
var appendPreReleaseCounter bool

var Command = cobra.Command{
	Use:   "next",
	Short: "prints version which should be used for the next release",
	Long:  "prints version which should be used for the next release according to commit message conventions and touched files",
	Run: func(cmd *cobra.Command, args []string) {

		nextVersion, err := next.Next(next.NextOptions{
			Workdir: common_opts.Workdir,
			Stable:  stable,
			PreReleaseOptions: semver.PreReleaseOptions{
				Label:         preReleaseTag,
				AppendCounter: appendPreReleaseCounter,
			},
		})

		if err != nil {
			logger.Logger.Fatalln(err)
		}

		fmt.Print(nextVersion.ToString())

	},
}

func init() {
	Command.Flags().BoolVar(&stable, "stable", true, "Specifies if this project is considered stable. Setting this to false will cause the major version to be 0. This command will fail if there is already a major version greater than 0.")
	Command.Flags().StringVar(&preReleaseTag, "pre-release-tag", "", "Specifies a pre-release tag which should be appended to the next version.")
	Command.Flags().BoolVar(&appendPreReleaseCounter, "pre-release-counter", false, "Specifies if there should be a counter appended to the pre-release tag. It will increase automatically depending on previous pre-releases for the same version.")
}
