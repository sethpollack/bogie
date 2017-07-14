package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"go.mozilla.org/sops/decrypt"
)

func proccessApplications(b *Bogie, outFile io.Writer) error {
	for _, app := range b.Applications {
		err := proccessApplication(app.Templates, app.Values, b.EnvFile, outFile, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func proccessApplication(input, values, envfile string, outFile io.Writer, b *Bogie) error {
	input = filepath.Clean(input)

	_, err := os.Stat(input)
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())

		if ok, _ := regexp.MatchString(b.IgnoreRegex, entry.Name()); ok {
			continue
		}

		if entry.IsDir() {
			err := proccessApplication(nextInPath, values, envfile, outFile, b)
			if err != nil {
				return err
			}
		} else {
			inString, err := readInput(nextInPath)
			if err != nil {
				return err
			}

			var inValues []byte
			if values != "" {
				inValues, err = decrypt.File(values, "yaml")
				if err != nil {
					return err
				}
			} else {
				log.Printf("No values found for template (%v)\n", nextInPath)
			}

			var inEnv []byte
			if envfile != "" {
				inEnv, err = decrypt.File(envfile, "yaml")
				if err != nil {
					return err
				}
			} else {
				log.Printf("No env_file found for template (%v)\n", nextInPath)
			}

			fmt.Fprint(outFile, "\n---\n")

			if err := renderTemplate(b, inString, string(inValues), string(inEnv), outFile); err != nil {
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
