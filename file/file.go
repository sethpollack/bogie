package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"go.mozilla.org/sops/decrypt"
)

func ReadFile(f func(text string, out io.Writer)) func(string) string {
	return func(path string) string {
		output, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		var buff bytes.Buffer
		f(string(output), &buff)
		return buff.String()
	}
}

func DecryptFile(f func(text string, out io.Writer)) func(string) string {
	return func(path string) string {
		output, err := decrypt.File(path, "yaml")
		if err != nil {
			log.Fatal(err)
		}

		var buff bytes.Buffer
		f(string(output), &buff)
		return buff.String()
	}
}

func ReadDir(f func(text string, out io.Writer)) func(string) map[string]string {
	readFileFunc := ReadFile(f)
	return func(dir string) map[string]string {
		fileMap := make(map[string]string)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if !file.IsDir() {
				fileMap[file.Name()] = readFileFunc(dir + "/" + file.Name())
			}
		}
		return fileMap
	}
}

func DecryptDir(f func(text string, out io.Writer)) func(string) map[string]string {
	readFileFunc := DecryptFile(f)
	return func(dir string) map[string]string {
		fileMap := make(map[string]string)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if !file.IsDir() {
				fileMap[file.Name()] = readFileFunc(dir + "/" + file.Name())
			}
		}
		return fileMap
	}
}
