package compare

import (
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
	"github.com/spf13/cobra"
)

var Command = cobra.Command{
	Use:   "compare <version-1> <version-2>",
	Short: "compares two semantic versions",
	Long: `This command is an utility command to compare two semantic versions. 

- Prints "=" if both versions are equal
- Prints "<" if the first version is less than the second version
- Prints ">" if the first version is greater than the second version
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			logger.Logger.Fatalln("Did not expect", len(args), "arguments")
		}

		v1, err := semver.ParseVersion(args[0])

		if err != nil {
			logger.Logger.Fatalln("Could not parse argument 1:", err)
		}

		v2, err := semver.ParseVersion(args[1])

		if err != nil {
			logger.Logger.Fatalln("Could not parse argument 2:", err)
		}

		result := semver.CompareVersions(v1, v2)

		switch result {
		case 0:
			print("=")
		case 1:
			print(">")
		case -1:
			print("<")
		}
	},
}
