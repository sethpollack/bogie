package main

import (
	"io"
	"log"
	"text/template"

	"github.com/Masterminds/sprig"
)

func (b *Bogie) createTemplate() *template.Template {
	return template.New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(b.funcMap).
		Option("missingkey=error")
}

type Bogie struct {
	funcMap    template.FuncMap
	leftDelim  string
	rightDelim string
}

func (b *Bogie) RunTemplate(text string, out io.Writer) {
	context := &Context{}
	tmpl, err := b.createTemplate().Delims(b.leftDelim, b.rightDelim).Parse(text)
	if err != nil {
		log.Fatalf("Line %q: %v\n", text, err)
	}
	if err := tmpl.Execute(out, context); err != nil {
		panic(err)
	}
}

func NewBogie(o *BogieOpts) *Bogie {
	return &Bogie{
		leftDelim:  o.lDelim,
		rightDelim: o.rDelim,
		funcMap:    initFuncs(o),
	}
}

func runTemplate(o *BogieOpts) error {
	b := NewBogie(o)
	if o.inputDir != "" {
		return processInputDir(o.inputDir, o.outputDir, b)
	}
	return processInputFiles(o.input, o.inputFiles, o.outputFiles, b)
}

func renderTemplate(b *Bogie, inString string, outPath string) error {
	outFile, err := openOutFile(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	b.RunTemplate(inString, outFile)
	return nil
}
