package bogie

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	dotaccess "github.com/go-bongo/go-dotaccess"
	"github.com/imdario/mergo"
	bogieio "github.com/sethpollack/bogie/io"
	yaml "gopkg.in/yaml.v2"
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
	re := regexp.MustCompile(b.AppRegex)
	for _, app := range b.ApplicationInputs {
		if b.AppRegex != "" {
			if !re.MatchString(app.Name) {
				continue
			}
		}

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

	inEnv, err := bogieio.DecryptFile(envfile, "yaml")
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

	dontWarn := func(file string) bool {
		return len(app.Values) == 0 &&
			file == fmt.Sprintf("%s/values.yaml", app.Templates)
	}

	for _, file := range files {
		b, err := bogieio.DecryptFile(file, "yaml")
		if err != nil {
			if dontWarn(file) {
				continue
			}
			return &c, err
		}

		var tmp map[interface{}]interface{}
		err = yaml.Unmarshal(b, &tmp)
		if err != nil {
			return &c, err
		}

		mergo.Merge(&c.Values, tmp)
	}

	mergo.Merge(&c.Values, old.Values)

	for _, keyVal := range app.OverrideVars {
		splits := strings.SplitN(keyVal, "=", 2)
		err := dotaccess.Set(c.Values, splits[0], splits[1])
		if err != nil {
			return &c, err
		}
	}

	return &c, nil
}

func processApplication(conf config) error {
	input := conf.input
	output := conf.output

	entries, err := bogieio.ReadDir(input)
	if err != nil {
		return err
	}

	helper, _ := bogieio.ReadFile(input + "/_helpers.tmpl")

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
			inString, err := bogieio.ReadFile(nextInPath)
			if err != nil {
				return err
			}

			*conf.appOutputs = append(*conf.appOutputs, &applicationOutput{
				outPath:  nextOutPath,
				template: string(helper) + string(inString),
				context:  conf.context,
			})
		}
	}

	return nil
}
