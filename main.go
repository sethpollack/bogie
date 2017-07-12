package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type BogieOpts struct {
	lDelim    string
	rDelim    string
	manifest  string
	templates string
	outDir    string
	envFile   string
	values    string
}

var opts BogieOpts

func validateOpts(cmd *cobra.Command, args []string) error {
	return nil
}

func newBogieCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "bogie",
		Short:   "Process text files with Go templates",
		PreRunE: validateOpts,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplate(&opts)
		},
	}
	return rootCmd
}

func initFlags(command *cobra.Command) {
	command.Flags().StringVar(&opts.lDelim, "left-delim", "{{{", "override the default left-`delimiter`")
	command.Flags().StringVar(&opts.rDelim, "right-delim", "}}}", "override the default right-`delimiter`")
	command.Flags().StringVar(&opts.manifest, "manifest", "", "template manifest")

	command.Flags().StringVar(&opts.envFile, "env-file", "", "global values file - required when not using a manifest.")
	command.Flags().StringVar(&opts.values, "values", "", "values file - required when not using a manifest.")
	command.Flags().StringVar(&opts.outDir, "output-dir", "releases", "`directory` to store the processed templates - required when not using a manifest.")
}

func main() {
	command := newBogieCmd()
	initFlags(command)
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
