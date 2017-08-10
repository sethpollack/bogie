package bogie

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/sethpollack/bogie/util"
	"go.mozilla.org/sops/decrypt"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	appOutputs  *[]*applicationOutput
	muteWarning bool
	input       string
	output      string
	c           *context
	b           *Bogie
}

func processApplications(b *Bogie) ([]*applicationOutput, error) {
	c, err := genContext(b.EnvFile)
	if err != nil {
		return nil, err
	}

	if c.Env == nil {
		log.Print("No env_file found")
	}

	appOutputs := []*applicationOutput{}

	for _, app := range b.ApplicationInputs {
		c, err := setValueContext(app.Values, c)
		if err != nil {
			return nil, err
		}

		releaseDir := filepath.Join(b.OutPath, app.Name)

		conf := Config{
			appOutputs:  &appOutputs,
			input:       app.Templates,
			output:      releaseDir,
			muteWarning: app.MuteWarning,
			c:           c,
			b:           b,
		}

		err = processApplication(conf)
		if err != nil {
			return nil, err
		}
	}

	return appOutputs, nil
}

func setValueContext(values string, c context) (*context, error) {
	if values != "" {
		inValues, err := decrypt.File(values, "yaml")
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal([]byte(inValues), &c.Values)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func genContext(envfile string) (context, error) {
	c := context{}

	if envfile != "" {
		inEnv, err := decrypt.File(envfile, "yaml")
		if err != nil {
			return context{}, err
		}

		err = yaml.Unmarshal([]byte(inEnv), &c.Env)
		if err != nil {
			return context{}, err
		}
	}

	return c, nil
}

func processApplication(conf Config) error {
	input := filepath.Clean(conf.input)
	output := filepath.Clean(conf.output)

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	helper, _ := util.ReadInput(input + "/_helpers.tmpl")

	r := conf.b.Rules.Clone()
	r.ParseFile(input + "/.bogieignore")

	for _, entry := range entries {
		if ok := r.Ignore(entry.Name(), entry.IsDir()); ok {
			continue
		}

		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(output, entry.Name())

		if entry.IsDir() {
			conf.input = nextInPath
			conf.output = nextOutPath
			err := processApplication(conf)
			if err != nil {
				return err
			}
		} else {
			inString, err := util.ReadInput(nextInPath)
			if err != nil {
				return err
			}

			if conf.c.Values == nil && !conf.muteWarning {
				log.Printf("No values found for template (%v)", nextInPath)
			}

			*conf.appOutputs = append(*conf.appOutputs, &applicationOutput{
				outPath:  nextOutPath,
				template: helper + inString,
				context:  conf.c,
			})
		}
	}

	return nil
}
