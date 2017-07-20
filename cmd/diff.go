package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Get diff between manifest and kube api",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("running diff")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
