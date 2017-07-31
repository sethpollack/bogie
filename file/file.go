package file

import (
	"bytes"
	"io"
	"io/ioutil"

	"go.mozilla.org/sops/decrypt"
)

func ReadFile(f func(text string, out io.Writer) error) func(string) (string, error) {
	return func(path string) (string, error) {
		output, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		var buff bytes.Buffer

		f(string(output), &buff)

		return buff.String(), nil
	}
}

func DecryptFile(f func(text string, out io.Writer) error) func(string) (string, error) {
	return func(path string) (string, error) {
		output, err := decrypt.File(path, "yaml")
		if err != nil {
			return "", err
		}

		var buff bytes.Buffer

		f(string(output), &buff)

		return buff.String(), nil
	}
}

func ReadDir(f func(text string, out io.Writer) error) func(string) (map[string]string, error) {
	readFileFunc := ReadFile(f)
	return func(dir string) (map[string]string, error) {
		fileMap := make(map[string]string)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !file.IsDir() {
				res, err := readFileFunc(dir + "/" + file.Name())
				if err != nil {
					return nil, err
				}

				fileMap[file.Name()] = res
			}
		}

		return fileMap, nil
	}
}

func DecryptDir(f func(text string, out io.Writer) error) func(string) (map[string]string, error) {
	readFileFunc := DecryptFile(f)
	return func(dir string) (map[string]string, error) {
		fileMap := make(map[string]string)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !file.IsDir() {
				res, err := readFileFunc(dir + "/" + file.Name())
				if err != nil {
					return nil, err
				}

				fileMap[file.Name()] = res
			}
		}

		return fileMap, nil
	}
}
