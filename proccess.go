package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"go.mozilla.org/sops/decrypt"
)

func proccessApplications(b *Bogie) ([]*ApplicationOutput, error) {
	appOutputs := []*ApplicationOutput{}

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

func setValueContext(values string, c *Context) (*Context, error) {
	nc := &Context{
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

func genContext(envfile string) (*Context, error) {
	c := &Context{}

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

func proccessApplication(input string, c *Context, b *Bogie) ([]*ApplicationOutput, error) {
	input = filepath.Clean(input)

	_, err := os.Stat(input)
	if err != nil {
		return nil, err
	}

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return nil, err
	}

	appOutputs := []*ApplicationOutput{}

	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(b.OutDir, input, entry.Name())

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
			inString, err := readInput(nextInPath)
			if err != nil {
				return nil, err
			}

			if c.Values == nil {
				log.Printf("No values found for template (%v)\n", nextInPath)
			}

			if c.Env == nil {
				log.Printf("No env_file found for template (%v)\n", nextInPath)
			}

			appOutputs = append(appOutputs, &ApplicationOutput{
				OutPath:  nextOutPath,
				Template: inString,
				Context:  c,
			})
		}
	}

	return appOutputs, nil
}

func readInput(filename string) (string, error) {
	var err error
	var inFile *os.File

	inFile, err = os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open %s\n%v", filename, err)
	}
	defer inFile.Close()

	bytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		err = fmt.Errorf("read failed for %s\n%v", filename, err)
		return "", err
	}

	return string(bytes), nil
}
