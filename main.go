package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type BogieOpts struct {
	lDelim      string
	rDelim      string
	inputIgnore string
	inputDir    string
	outputDir   string
	envfile     string
}

var opts BogieOpts

func validateOpts(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("env-file").Changed {
		return errors.New("--env-file must be set")
	}

	if !cmd.Flag("input-dir").Changed {
		return errors.New("--input-dir must be set")
	}

	if !cmd.Flag("output-dir").Changed {
		return errors.New("--output-dir must be set")
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
	command.Flags().StringVar(&opts.envfile, "env-file", "", "app env file to load into the context must live at the base of the templates dir")
	command.Flags().StringVar(&opts.inputDir, "input-dir", "templates", "`directory` which is examined recursively for templates")
	command.Flags().StringVar(&opts.outputDir, "output-dir", "releases", "`directory` to store the processed templates.")
	command.Flags().StringVar(&opts.inputIgnore, "input-ignore", "values.yaml", "regex for files/directories to ignore")

	env := &Env{}
	ldDefault := env.Getenv("BOGIE_LEFT_DELIM", "{{{")
	rdDefault := env.Getenv("BOGIE_RIGHT_DELIM", "}}}")

	command.Flags().StringVar(&opts.lDelim, "left-delim", ldDefault, "override the default left-`delimiter` [$BOGIE_LEFT_DELIM]")
	command.Flags().StringVar(&opts.rDelim, "right-delim", rdDefault, "override the default right-`delimiter` [$BOGIE_RIGHT_DELIM]")
}

func main() {
	command := newBogieCmd()
	initFlags(command)
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
