package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var kubeconfig string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Get diff between manifest and kube api",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running diff.")
	},
}

func init() {
	templateCmd.AddCommand(diffCmd)
	templateCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "/root/.kube/config", "path to kube/config")
}
