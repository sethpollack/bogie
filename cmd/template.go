package cmd

import (
	"errors"
	"fmt"

	"github.com/sethpollack/bogie/bogie"
	"github.com/sethpollack/bogie/io"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type bogieOpts struct {
	lDelim          string
	rDelim          string
	environment     string
	manifest        string
	outFormat       string
	templates       string
	outPath         string
	outFile         string
	envFile         string
	valuesFile      []string
	overrideVars    []string
	templatesPath   string
	ignoreFile      string
	skipImageLookup bool
}

var o bogieOpts

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVar(&o.lDelim, "left-delim", "{{{", "override the default left-`delimiter`")
	templateCmd.Flags().StringVar(&o.rDelim, "right-delim", "}}}", "override the default right-`delimiter`")

	templateCmd.Flags().StringVar(&o.environment, "env", "", "environment")
	templateCmd.Flags().StringVarP(&o.manifest, "manifest", "m", "", "template manifest")
	templateCmd.Flags().StringVarP(&o.outFormat, "out", "o", "stdout", "output format (dir|file|stdout)")

	templateCmd.Flags().StringVarP(&o.envFile, "env-file", "e", "", "global values file.")
	templateCmd.Flags().StringSliceVarP(&o.valuesFile, "values-file", "v", []string{}, "values file.")
	templateCmd.Flags().StringSliceVarP(&o.overrideVars, "override-vars", "x", []string{}, "extra vars path.to.key=value format.")
	templateCmd.Flags().StringVarP(&o.templatesPath, "templates-path", "t", "", "templates.")

	templateCmd.Flags().StringVarP(&o.outPath, "output-path", "p", "releases", "`dir` to store the processed templates.")
	templateCmd.Flags().StringVarP(&o.outFile, "output-file", "f", "release.yaml", "`file` to store the processed templates.")

	templateCmd.Flags().StringVarP(&o.ignoreFile, "ignore-file", "i", ".bogieignore", ".bogieignore file")
	templateCmd.Flags().BoolVarP(&o.skipImageLookup, "skip-image-lookup", "s", false, "Skip image lookup in template function latestImage")
}

var template_example = `
# example single run with var files
bogie template \
 -t path/to/templates \
 -v path/to/templates/env.values.yaml \
 -v path/to/templates/values.yaml \
 -e path/to/global/vars/values.yaml \
 -o dir

# example single run with auto var files
bogie template \
 -t path/to/templates \
 -e path/to/global/vars/values.yaml \
 --env production \
 -o dir

# example manifest run
bogie template \
 -m path/to/manifest.yaml

#example manifest
out_path: releases
out_file: release.yaml
out_format: dir
env_file: path/to/global/vars/values.yaml
ignore_file: .bogieignore
skip_image_lookup: false
applications:
- name: my-templates
  templates: path/to/templates
  values:
  - path/to/templates/values.yaml
  - path/to/templates/env.values.yaml
  override_vars:
  - app.secrets.key=value
- name: my-other-templates
  templates: path/to/templates
  env: production
  override_vars:
  - app.secrets.key=value
- name: my-other-other-templates
  templates: path/to/templates
  mute_warning: true
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
		b, err := newBogie(&o)
		if err != nil {
			return err
		}

		return b.Run()
	},
}

func newBogie(o *bogieOpts) (*bogie.Bogie, error) {
	b := &bogie.Bogie{
		EnvFile:         o.envFile,
		OutPath:         o.outPath,
		OutFile:         o.outFile,
		OutFormat:       o.outFormat,
		LDelim:          o.lDelim,
		RDelim:          o.rDelim,
		SkipImageLookup: o.skipImageLookup,
	}

	b.InitRules()
	b.Rules.ParseFile(o.ignoreFile)

	if o.templatesPath != "" {
		b.ApplicationInputs = []*bogie.ApplicationInput{
			{
				Templates:    o.templatesPath,
				Values:       o.valuesFile,
				Env:          o.environment,
				OverrideVars: o.overrideVars,
			},
		}
	}

	if o.manifest != "" {
		err := parseManifest(o.manifest, b)
		if err != nil {
			return &bogie.Bogie{}, errors.New(fmt.Sprintf("error parsing manifest file %v\n", err))
		}
	}

	return b, nil
}

func parseManifest(manifest string, b *bogie.Bogie) error {
	output, err := io.ReadInput(manifest)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(output), b)
	if err != nil {
		return err
	}

	return nil
}
