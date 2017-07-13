package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"go.mozilla.org/sops/decrypt"
)

func proccessApplications(b *Bogie) error {
	for _, app := range b.Applications {
		err := proccessApplication(app.Templates, b.OutDir+"/"+app.Name, app.Values, b.EnvFile, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func proccessApplication(input, output, values, envfile string, b *Bogie) error {
	input = filepath.Clean(input)
	output = filepath.Clean(output)

	fi, err := os.Stat(input)
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(output, fi.Mode()); err != nil {
		return err
	}

	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(output, entry.Name())

		if ok, _ := regexp.MatchString(b.IgnoreRegex, entry.Name()); ok {
			continue
		}

		if entry.IsDir() {
			err := proccessApplication(nextInPath, nextOutPath, values, envfile, b)
			if err != nil {
				return err
			}
		} else {
			inString, err := readInput(nextInPath)
			if err != nil {
				return err
			}

			inValues, err := decrypt.File(values, "yaml")
			if err != nil {
				return err
			}

			inEnv, err := decrypt.File(envfile, "yaml")
			if err != nil {
				return err
			}

			if err := renderTemplate(b, inString, string(inValues), string(inEnv), nextOutPath); err != nil {
				return err
			}
		}
	}
	return nil
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

func openOutFile(filename string) (out *os.File, err error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}
