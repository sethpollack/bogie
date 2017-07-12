package main

import (
	"io"
	"log"
	"text/template"

	"github.com/Masterminds/sprig"
	yaml "gopkg.in/yaml.v2"
)

var temp *template.Template

func (b *Bogie) createTemplate() *template.Template {
	if temp == nil {
		temp = template.New("template").
			Funcs(sprig.TxtFuncMap()).
			Funcs(b.InitFuncs()).
			Option("missingkey=error")
	}
	return temp
}

type Bogie struct {
	leftDelim   string
	rightDelim  string
	inputIgnore string
	inputDir    string
	context     *Context
}

type Context struct {
	Values *map[interface{}]interface{}
	Env    *map[interface{}]interface{}
}

func (b *Bogie) RunTemplate(text string, out io.Writer) {
	tmpl, err := b.createTemplate().Delims(b.leftDelim, b.rightDelim).Parse(text)
	if err != nil {
		log.Fatalf("Line %q: %v\n", text, err)
	}
	if err := tmpl.Execute(out, b.context); err != nil {
		panic(err)
	}
}

func NewBogie(o *BogieOpts) *Bogie {
	return &Bogie{
		leftDelim:   o.lDelim,
		rightDelim:  o.rDelim,
		inputDir:    o.inputDir,
		inputIgnore: o.inputIgnore,
		context:     &Context{},
	}
}

func runTemplate(o *BogieOpts) error {
	b := NewBogie(o)
	return processInputDir(o.envfile, o.inputDir, o.outputDir, b)
}

func renderTemplate(b *Bogie, inString, inValues, inEnv, outPath string) error {
	outFile, err := openOutFile(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = yaml.Unmarshal([]byte(inValues), &b.context.Values)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(inEnv), &b.context.Env)
	if err != nil {
		return err
	}

	b.RunTemplate(inString, outFile)
	return nil
}
