package next

import (
	"fmt"
	"github.com/psanetra/git-semver/cli/common_opts"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/next"
	"github.com/psanetra/git-semver/semver"
	"github.com/spf13/cobra"
)

var stable bool
var majorVersionFilter int
var preReleaseTag string
var appendPreReleaseCounter bool

var Command = cobra.Command{
	Use:   "next",
	Short: "prints version which should be used for the next release",
	Long:  `This command can be used to calculate the next semantic version based on the history of the current branch. It fails if the git tag of the latest semantic version is not reachable on the current branch or if the tagged commit is not reachable because the repository is shallow.`,
	Run: func(cmd *cobra.Command, args []string) {

		nextVersion, err := next.Next(next.NextOptions{
			Workdir: common_opts.Workdir,
			Stable:  stable,
			MajorVersionFilter: majorVersionFilter,
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
	Command.Flags().IntVar(&majorVersionFilter, "major-version", -1, "Only consider tags with this specific major version.")
	Command.Flags().StringVar(&preReleaseTag, "pre-release-tag", "", "Specifies a pre-release tag which should be appended to the next version.")
	Command.Flags().BoolVar(&appendPreReleaseCounter, "pre-release-counter", false, "Specifies if there should be a counter appended to the pre-release tag. It will increase automatically depending on previous pre-releases for the same version.")
}
