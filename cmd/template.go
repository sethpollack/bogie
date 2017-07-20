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
	templateCmd.Flags().StringVarP(&o.TemplatesPath, "templates-path", "t", "", "templates.")

	templateCmd.Flags().StringVarP(&o.OutPath, "output-path", "p", "releases", "`dir` to store the processed templates.")
	templateCmd.Flags().StringVarP(&o.OutFile, "output-file", "f", "release.yaml", "`file` to store the processed templates.")

	templateCmd.Flags().StringVarP(&o.IgnoreRegex, "ignore-regex", "i", "((.+).md|(.+)?values.yaml)", "regex to skip files from being copied over.")
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Process text files with Go templates",
	Long: `
	# example single run
		bogie template \
			-t path/to/templates \
			-v path/to/templates/values.yaml \
			-e path/to/global/vars/values.yaml \
			-o file

		# example manifest run
		bogie template \
			-m path/to/manifest.yaml \
			-o file

		#example manifest
		out_path: releases
		out_file: release.yaml
		env_file: path/to/global/vars/values.yaml
		ignore_regex: (.+)?values.yaml
		applications:
		- name: my-templates
  		templates: path/to/templates
			values: path/to/templates/values.yaml
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if o.Manifest == "" {
			if o.TemplatesPath == "" {
				return errors.New("--templates-path is required when not using the manifest file")
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return bogie.RunBogie(&o)
	},
}
