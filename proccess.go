package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func processInputFiles(stringTemplate string, input []string, output []string, b *Bogie) error {
	input, err := readInputs(stringTemplate, input)
	if err != nil {
		return err
	}

	if len(output) == 0 {
		output = []string{"-"}
	}

	for n, input := range input {
		if err := renderTemplate(b, input, output[n]); err != nil {
			return err
		}
	}
	return nil
}

func processInputDir(input string, output string, b *Bogie) error {
	input = filepath.Clean(input)
	output = filepath.Clean(output)

	si, err := os.Stat(input)
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(output, si.Mode()); err != nil {
		return err
	}

	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(output, entry.Name())

		if entry.IsDir() {
			err := processInputDir(nextInPath, nextOutPath, b)
			if err != nil {
				return err
			}
		} else {
			inString, err := readInput(nextInPath)
			if err != nil {
				return err
			}
			if err := renderTemplate(b, inString, nextOutPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func readInputs(input string, files []string) ([]string, error) {
	if input != "" {
		return []string{input}, nil
	}
	if len(files) == 0 {
		files = []string{"-"}
	}
	ins := make([]string, len(files))

	for n, filename := range files {
		inString, err := readInput(filename)
		if err != nil {
			return nil, err
		}
		ins[n] = inString
	}
	return ins, nil
}

func readInput(filename string) (string, error) {
	var err error
	var inFile *os.File
	if filename == "-" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("failed to open %s\n%v", filename, err)
		}
		defer inFile.Close()
	}
	bytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		err = fmt.Errorf("read failed for %s\n%v", filename, err)
		return "", err
	}
	return string(bytes), nil
}

func openOutFile(filename string) (out *os.File, err error) {
	if filename == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}
