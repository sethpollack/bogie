package bogie

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/sethpollack/bogie/util"

	yaml "gopkg.in/yaml.v2"

	"go.mozilla.org/sops/decrypt"
)

func proccessApplications(b *Bogie) ([]*applicationOutput, error) {
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

		err = proccessApplication(&appOutputs, app.Templates, c, b)
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

func proccessApplication(appOutputs *[]*applicationOutput, input string, c *context, b *Bogie) error {
	input = filepath.Clean(input)

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(b.OutPath, input, entry.Name())

		if ok, _ := regexp.MatchString(b.IgnoreRegex, entry.Name()); ok {
			continue
		}

		if entry.IsDir() {
			err := proccessApplication(appOutputs, nextInPath, c, b)
			if err != nil {
				return err
			}
		} else {
			inString, err := util.ReadInput(nextInPath)
			if err != nil {
				return err
			}

			if c.Values == nil {
				log.Printf("No values found for template (%v)", nextInPath)
			}

			*appOutputs = append(*appOutputs, &applicationOutput{
				outPath:  nextOutPath,
				template: inString,
				context:  c,
			})
		}
	}

	return nil
}
