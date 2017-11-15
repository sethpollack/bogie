package cmd

import (
	"fmt"

	"github.com/sethpollack/bogie/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s\nCommit %s\n", version.Version, version.Commit)
	},
}
