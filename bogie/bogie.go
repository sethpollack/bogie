package bogie

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
)

type ApplicationInput struct {
	Templates string
	Values    string
}

type applicationOutput struct {
	outPath  string
	template string
	context  *context
}

type context struct {
	Values map[interface{}]interface{}
	Env    map[interface{}]interface{}
}

type Bogie struct {
	RDelim            string
	LDelim            string
	EnvFile           string              `yaml:"env_file"`
	OutFile           string              `yaml:"out_file"`
	OutPath           string              `yaml:"out_path"`
	OutFormat         string              `yaml:"out_format"`
	IgnoreRegex       string              `yaml:"ignore_regex"`
	ApplicationInputs []*ApplicationInput `yaml:"applications"`
}

func RunBogie(b *Bogie) error {
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

func runTemplate(c *context, b *Bogie, text string, out io.Writer) {
	tmpl, err := template.New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(initFuncs(c, b)).
		Option("missingkey=error").
		Delims(b.LDelim, b.RDelim).
		Parse(text)

	if err != nil {
		log.Fatalf("Line %q: %v\n", text, err)
	}

	if err := tmpl.Execute(out, c); err != nil {
		panic(err)
	}
}

func renderTemplateToDir(b *Bogie, apps []*applicationOutput) error {
	for _, app := range apps {
		if err := os.MkdirAll(path.Dir(app.outPath), os.FileMode(0744)); err != nil {
			return err
		}

		f, err := openOutFile(app.outPath)
		if err != nil {
			return err
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		defer w.Flush()

		runTemplate(app.context, b, app.template, w)
	}

	return nil
}

func renderTemplateToFile(b *Bogie, apps []*applicationOutput) error {
	if err := os.MkdirAll(b.OutPath, os.FileMode(0744)); err != nil {
		return err
	}

	f, err := openOutFile(fmt.Sprintf("%s/%s", b.OutPath, b.OutFile))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, app := range apps {
		fmt.Fprint(w, "\n---\n")
		runTemplate(app.context, b, app.template, w)
	}

	return nil
}

func renderTemplateToSTDOUT(b *Bogie, apps []*applicationOutput) error {
	out := os.Stdout
	for _, app := range apps {
		fmt.Fprint(out, "\n---\n")
		runTemplate(app.context, b, app.template, out)
	}

	return nil
}

func openOutFile(filename string) (out *os.File, err error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}
