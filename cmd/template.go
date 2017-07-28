package cmd

import (
	"errors"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/sethpollack/bogie/bogie"
	"github.com/sethpollack/bogie/ignore"
	"github.com/sethpollack/bogie/util"
	"github.com/spf13/cobra"
)

type bogieOpts struct {
	lDelim        string
	rDelim        string
	manifest      string
	outFormat     string
	templates     string
	outPath       string
	outFile       string
	envFile       string
	valuesFile    string
	templatesPath string
	ignoreFile    string
}

var o bogieOpts

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.PersistentFlags().StringVar(&o.lDelim, "left-delim", "{{{", "override the default left-`delimiter`")
	templateCmd.PersistentFlags().StringVar(&o.rDelim, "right-delim", "}}}", "override the default right-`delimiter`")

	templateCmd.PersistentFlags().StringVarP(&o.manifest, "manifest", "m", "", "template manifest")
	templateCmd.PersistentFlags().StringVarP(&o.outFormat, "out", "o", "dir", "output format (dir|file|stdout)")

	templateCmd.PersistentFlags().StringVarP(&o.envFile, "env-file", "e", "", "global values file.")
	templateCmd.PersistentFlags().StringVarP(&o.valuesFile, "values-file", "v", "", "values file.")
	templateCmd.PersistentFlags().StringVarP(&o.templatesPath, "templates-path", "t", "", "templates.")

	templateCmd.PersistentFlags().StringVarP(&o.outPath, "output-path", "p", "releases", "`dir` to store the processed templates.")
	templateCmd.PersistentFlags().StringVarP(&o.outFile, "output-file", "f", "release.yaml", "`file` to store the processed templates.")

	templateCmd.PersistentFlags().StringVarP(&o.ignoreFile, "ignore-file", "i", ".bogieignore", ".bogieignore file")
}

var template_example = `
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
ignore_file: .bogieignore
applications:
- name: my-templates
  templates: path/to/templates
  values: path/to/templates/values.yaml

`

var templateCmd = &cobra.Command{
	Use:     "template",
	Short:   "Process text files with Go templates",
	Example: template_example,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if o.manifest == "" {
			if o.templatesPath == "" {
				return errors.New("--templates-path is required when not using the manifest file")
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		b := newBogie(&o)
		return bogie.RunBogie(b)
	},
}

func newBogie(o *bogieOpts) *bogie.Bogie {
	b := &bogie.Bogie{
		EnvFile:   o.envFile,
		OutPath:   o.outPath,
		OutFile:   o.outFile,
		OutFormat: o.outFormat,
		LDelim:    o.lDelim,
		RDelim:    o.rDelim,
	}

	r := ignore.Init()
	r.ParseFile(o.ignoreFile)
	b.Rules = r

	if o.templatesPath != "" && o.valuesFile != "" {
		b.ApplicationInputs = []*bogie.ApplicationInput{
			{Templates: o.templatesPath, Values: o.valuesFile},
		}
	}

	if o.manifest != "" {
		err := parseManifest(o.manifest, b)
		if err != nil {
			log.Fatalf("error parsing manifest file %v\n", err)
		}
	}

	return b
}

func parseManifest(manifest string, b *bogie.Bogie) error {
	output, err := util.ReadInput(manifest)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(output), b)
	if err != nil {
		return err
	}

	return nil
}
