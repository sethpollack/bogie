package cmd

import (
	"errors"

	"github.com/sethpollack/bogie/bogie"
	"github.com/spf13/cobra"
)

var o bogie.BogieOpts

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVar(&o.LDelim, "left-delim", "{{{", "override the default left-`delimiter`")
	templateCmd.Flags().StringVar(&o.RDelim, "right-delim", "}}}", "override the default right-`delimiter`")

	templateCmd.Flags().StringVarP(&o.Manifest, "manifest", "m", "", "template manifest")
	templateCmd.Flags().StringVarP(&o.OutFormat, "out", "o", "dir", "output format")

	templateCmd.Flags().StringVarP(&o.EnvFile, "env-file", "e", "", "global values file.")
	templateCmd.Flags().StringVarP(&o.ValuesFile, "values-file", "v", "", "values file.")
	templateCmd.Flags().StringVarP(&o.TemplatesPath, "templates-dir", "t", "", "templates dir.")

	templateCmd.Flags().StringVarP(&o.OutPath, "output-path", "p", "releases", "`dir` to store the processed templates.")
	templateCmd.Flags().StringVarP(&o.OutFile, "output-file", "f", "release.yaml", "`file` to store the processed templates.")

	templateCmd.Flags().StringVarP(&o.IgnoreRegex, "ignore-regex", "i", "((.+).md|(.+)?values.yaml)", "regex to skip files from being copied over.")
}

var templateCmd = &cobra.Command{
	Use:     "template",
	Short:   "Process text files with Go templates",
	Long:    ``,
	PreRunE: validateOpts,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bogie.RunBogie(&o)
	},
}

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
