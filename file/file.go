package file

import (
	"bytes"
	"io"

	bogieio "github.com/sethpollack/bogie/io"
)

var templater func(text string, out io.Writer) error

func SetTemplater(f func(text string, out io.Writer) error) {
	templater = f
}

func readFile(read func() ([]byte, error)) (string, error) {
	output, err := read()
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer

	err = templater(string(output), &buff)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func readDir(dir string, read func(string) (string, error)) (map[string]string, error) {
	fileMap := make(map[string]string)
	files, err := bogieio.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			res, err := read(dir + "/" + file.Name())

			if err != nil {
				return nil, err
			}

			fileMap[file.Name()] = res
		}
	}

	return fileMap, nil
}

func ReadFile(path string) (string, error) {
	return readFile(func() ([]byte, error) {
		return bogieio.ReadFile(path)
	})
}

func DecryptFile(path string) (string, error) {
	return readFile(func() ([]byte, error) {
		return bogieio.DecryptFile(path, "yaml")
	})
}

func ReadDir(dir string) (map[string]string, error) {
	return readDir(dir, ReadFile)
}

func DecryptDir(dir string) (map[string]string, error) {
	return readDir(dir, DecryptFile)
}
