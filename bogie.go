package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
	yaml "gopkg.in/yaml.v2"
)

var temp *template.Template

type ApplicationInput struct {
	Name      string
	Templates string
	Values    string
}

type ApplicationOutput struct {
	OutPath  string
	Template string
	Context  *Context
}

type Context struct {
	Values *map[interface{}]interface{}
	Env    *map[interface{}]interface{}
}

type Bogie struct {
	EnvFile           string `yaml:"env_file"`
	OutFile           string `yaml:"out_file"`
	OutDir            string `yaml:"out_dir"`
	OutFormat         string `yaml: "out_format"`
	LDelim            string
	RDelim            string
	IgnoreRegex       string              `yaml:"ignore_regex"`
	ApplicationInputs []*ApplicationInput `yaml:"applications"`
}

func NewBogie(o *BogieOpts) *Bogie {
	b := &Bogie{
		EnvFile:     o.envFile,
		OutDir:      o.outPath,
		OutFile:     o.outFile,
		OutFormat:   o.outFormat,
		LDelim:      o.lDelim,
		RDelim:      o.rDelim,
		IgnoreRegex: o.ignoreRegex,
	}

	if o.templatesPath != "" && o.valuesFile != "" {
		b.ApplicationInputs = append(b.ApplicationInputs, &ApplicationInput{
			Templates: o.templatesPath,
			Values:    o.valuesFile,
		})
	}

	if o.manifest != "" {
		err := parseManifest(o.manifest, b)
		if err != nil {
			log.Fatalf("error parsing manifest file %v\n", err)
		}
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
	apps, err := proccessApplications(b)
	if err != nil {
		return err
	}

	switch b.OutFormat {
	case "dir":
		return renderTemplateToDir(b, apps)
	case "file":
		return renderTemplateToFile(b, apps)
	default:
		return renderTemplateToSTDOUT(b, apps)
	}
}

func renderTemplateToDir(b *Bogie, apps []*ApplicationOutput) error {
	for _, app := range apps {

		if err := os.MkdirAll(path.Dir(app.OutPath), os.FileMode(0777)); err != nil {
			return err
		}

		f, err := openOutFile(app.OutPath)
		if err != nil {
			return err
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		defer w.Flush()

		RunTemplate(app.Context, b, app.Template, w)
	}

	return nil
}

func renderTemplateToFile(b *Bogie, apps []*ApplicationOutput) error {
	if err := os.MkdirAll(b.OutDir, os.FileMode(0777)); err != nil {
		return err
	}

	f, err := openOutFile(fmt.Sprintf("%s/%s", b.OutDir, b.OutFile))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, app := range apps {
		fmt.Fprint(w, "---\n")
		RunTemplate(app.Context, b, app.Template, w)
	}

	return nil
}

func renderTemplateToSTDOUT(b *Bogie, apps []*ApplicationOutput) error {
	for _, app := range apps {
		RunTemplate(app.Context, b, app.Template, os.Stdout)
	}

	return nil
}

func openOutFile(filename string) (out *os.File, err error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}
