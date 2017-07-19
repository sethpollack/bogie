package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type BogieOpts struct {
	lDelim string
	rDelim string

	manifest  string
	outFormat string

	//flags when not using manifest
	templates     string
	outPath       string
	outFile       string
	envFile       string
	valuesFile    string
	templatesPath string
	ignoreRegex   string
}

var opts BogieOpts

func validateOpts(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("manifest").Changed {
		if !cmd.Flag("output-file").Changed {
			return errors.New("--output-file is required when not using the manifest file")
		}

		if !cmd.Flag("templates-dir").Changed {
			return errors.New("--templates-dir is required when not using the manifest file")
		}
	}
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

	command.Flags().StringVar(&opts.outFormat, "out", "dir", "output format")

	command.Flags().StringVar(&opts.outPath, "output-dir", "releases", "`dir` to store the processed templates - required when not using a manifest.")
	command.Flags().StringVar(&opts.outFile, "output-file", "release.yaml", "`file` to store the processed templates - required when not using a manifest.")

	command.Flags().StringVar(&opts.templatesPath, "templates-dir", "", "templates dir - required when not using a manifest.")

	command.Flags().StringVar(&opts.envFile, "env-file", "", "global values file - used when not using a manifest (optional).")
	command.Flags().StringVar(&opts.valuesFile, "values-file", "", "values file - used when not using a manifest (optional).")

	command.Flags().StringVar(&opts.ignoreRegex, "ignore-regex", "((.+).md|(.+)?values.yaml)", "regex to skip files from being copied over - used when not using a manifest (optional).")
}

func main() {
	command := newBogieCmd()
	initFlags(command)
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
