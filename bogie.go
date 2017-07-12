package main

import (
	"io"
	"log"
	"text/template"

	"github.com/Masterminds/sprig"
	yaml "gopkg.in/yaml.v2"
)

var temp *template.Template

type Application struct {
	Name      string
	Templates string
	Values    string
}

type Context struct {
	Values *map[interface{}]interface{}
	Env    *map[interface{}]interface{}
}

type Bogie struct {
	EnvFile      string `yaml:"env_file"`
	OutDir       string `yaml:"out_dir"`
	LDelim       string
	RDelim       string
	IgnoreRegex  string `yaml:"ignore_regex"`
	Applications []*Application
}

func NewBogie(o *BogieOpts) *Bogie {
	b := &Bogie{
		EnvFile: o.envFile,
		OutDir:  o.outDir,
		LDelim:  o.lDelim,
		RDelim:  o.rDelim,
	}

	if o.templates != "" && o.values != "" {
		b.Applications = append(b.Applications, &Application{
			Templates: o.templates,
			Values:    o.values,
		})
	}

	err := parseManifest(o.manifest, b)
	if err != nil {
		log.Fatalf("error parsing manifest file %v\n", err)
	}
	return b
}

func createTemplate(c *Context, b *Bogie) *template.Template {
	if temp == nil {
		temp = template.New("template").
			Funcs(sprig.TxtFuncMap()).
			Funcs(initFuncs(c, b)).
			Option("missingkey=error")
	}
	return temp
}

func RunTemplate(c *Context, b *Bogie, text string, out io.Writer) {
	tmpl, err := createTemplate(c, b).
		Delims(b.LDelim, b.RDelim).
		Parse(text)

	if err != nil {
		log.Fatalf("Line %q: %v\n", text, err)
	}
	if err := tmpl.Execute(out, c); err != nil {
		panic(err)
	}
}

func parseManifest(manifest string, b *Bogie) error {
	output, err := readInput(manifest)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(output), b)
	if err != nil {
		return err
	}

	return nil
}

func runTemplate(o *BogieOpts) error {
	b := NewBogie(o)
	return proccessApplications(b)
}

func renderTemplate(b *Bogie, inString, inValues, inEnv, outPath string) error {
	outFile, err := openOutFile(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	c := &Context{}

	err = yaml.Unmarshal([]byte(inValues), &c.Values)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(inEnv), &c.Env)
	if err != nil {
		return err
	}

	RunTemplate(c, b, inString, outFile)
	return nil
}
