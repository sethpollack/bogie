package bogie

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"

	dotaccess "github.com/go-bongo/go-dotaccess"
	"github.com/imdario/mergo"
	"github.com/sethpollack/bogie/io"
)

type applicationOutput struct {
	outPath  string
	template string
	context  *context
}

type context struct {
	Values map[interface{}]interface{}
}

type config struct {
	appOutputs *[]*applicationOutput
	input      string
	output     string
	context    *context
	bogie      *Bogie
}

func processApplications(b *Bogie) ([]*applicationOutput, error) {
	c, err := genContext(b.EnvFile)
	if err != nil {
		return nil, err
	}

	appOutputs := []*applicationOutput{}
	for _, app := range b.ApplicationInputs {
		c, err := setValueContext(app, c)
		if err != nil {
			return nil, err
		}

		releaseDir := filepath.Join(b.OutPath, app.Name)

		conf := config{
			appOutputs: &appOutputs,
			input:      app.Templates,
			output:     releaseDir,
			context:    c,
			bogie:      b,
		}

		err = processApplication(conf)
		if err != nil {
			return nil, err
		}
	}

	return appOutputs, nil
}

func genContext(envfile string) (*context, error) {
	c := context{}

	if envfile == "" {
		return &c, nil
	}

	inEnv, err := io.DecryptFile(envfile, "yaml")
	if err != nil {
		return &c, err
	}

	err = yaml.Unmarshal(inEnv, &c.Values)
	if err != nil {
		return &c, err
	}

	return &c, nil
}

func setValueContext(app *ApplicationInput, old *context) (*context, error) {
	c := context{}

	files := []string{}

	if app.Env != "" {
		files = append(files, fmt.Sprintf("%s/%s.values.yaml", app.Templates, app.Env))
	}

	if len(app.Values) == 0 {
		files = append(files, fmt.Sprintf("%s/values.yaml", app.Templates))
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(app.Values)))
		files = append(files, app.Values...)
	}

	for _, file := range files {
		b, err := io.DecryptFile(file, "yaml")
		if err != nil {
			continue
		}
		var tmp map[interface{}]interface{}
		err = yaml.Unmarshal(b, &tmp)
		if err != nil {
			continue
		}

		mergo.Merge(&c.Values, tmp)
	}

	mergo.Merge(&c.Values, old.Values)

	for _, keyVal := range app.OverrideVars {
		splits := strings.SplitN(keyVal, "=", 2)
		err := dotaccess.Set(c.Values, splits[0], splits[1])
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func processApplication(conf config) error {
	input := conf.input
	output := conf.output

	entries, err := io.ReadDir(input)
	if err != nil {
		return err
	}

	helper, _ := io.ReadInput(input + "/_helpers.tmpl")

	r := conf.bogie.Rules.Clone()
	r.ParseFile(input + "/.bogieignore")

	for _, entry := range entries {
		if ok := r.Ignore(entry.Name(), entry.IsDir()); ok {
			continue
		}

		nextInPath := fmt.Sprintf("%s/%s", input, entry.Name())
		nextOutPath := filepath.Join(output, entry.Name())

		if entry.IsDir() {
			conf.input = nextInPath
			conf.output = nextOutPath

			err := processApplication(conf)
			if err != nil {
				return err
			}
		} else {
			inString, err := io.ReadInput(nextInPath)
			if err != nil {
				return err
			}

			*conf.appOutputs = append(*conf.appOutputs, &applicationOutput{
				outPath:  nextOutPath,
				template: helper + inString,
				context:  conf.context,
			})
		}
	}

	return nil
}
