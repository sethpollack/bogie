package bogie

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sethpollack/bogie/util"

	yaml "gopkg.in/yaml.v2"

	"go.mozilla.org/sops/decrypt"
)

func proccessApplications(b *Bogie) ([]*applicationOutput, error) {
	appOutputs := []*applicationOutput{}

	c, err := genContext(b.EnvFile)
	if err != nil {
		return nil, err
	}

	for _, app := range b.ApplicationInputs {
		c, err := setValueContext(app.Values, c)
		if err != nil {
			return nil, err
		}

		apps, err := proccessApplication(app.Templates, c, b)
		if err != nil {
			return nil, err
		}

		appOutputs = append(appOutputs, apps...)
	}

	return appOutputs, nil
}

func setValueContext(values string, c *context) (*context, error) {
	nc := &context{
		Env: c.Env,
	}

	if values != "" {
		inValues, err := decrypt.File(values, "yaml")
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal([]byte(inValues), &nc.Values)
		if err != nil {
			return nil, err
		}
	}

	return nc, nil
}

func genContext(envfile string) (*context, error) {
	c := &context{}

	if envfile != "" {
		inEnv, err := decrypt.File(envfile, "yaml")
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal([]byte(inEnv), &c.Env)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func proccessApplication(input string, c *context, b *Bogie) ([]*applicationOutput, error) {
	input = filepath.Clean(input)

	_, err := os.Stat(input)
	if err != nil {
		return nil, err
	}

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return nil, err
	}

	appOutputs := []*applicationOutput{}

	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(b.OutPath, input, entry.Name())

		if ok, _ := regexp.MatchString(b.IgnoreRegex, entry.Name()); ok {
			continue
		}

		if entry.IsDir() {
			apps, err := proccessApplication(nextInPath, c, b)
			if err != nil {
				return nil, err
			}

			appOutputs = append(appOutputs, apps...)
		} else {
			inString, err := util.ReadInput(nextInPath)
			if err != nil {
				return nil, err
			}

			if c.Values == nil {
				log.Printf("No values found for template (%v)\n", nextInPath)
			}

			if c.Env == nil {
				log.Printf("No env_file found for template (%v)\n", nextInPath)
			}

			appOutputs = append(appOutputs, &applicationOutput{
				outPath:  nextOutPath,
				template: inString,
				context:  c,
			})
		}
	}

	return appOutputs, nil
}
