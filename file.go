package main

import (
	"bytes"
	"io/ioutil"
	"log"
)

type File struct {
}

func (f *File) ReadFile(b *Bogie) func(string) string {
	return func(path string) string {
		output, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		var buff bytes.Buffer
		b.RunTemplate(string(output), &buff)
		return buff.String()
	}
}

func (f *File) ReadDir(b *Bogie) func(string) map[string]string {
	readFileFunc := f.ReadFile(b)
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
