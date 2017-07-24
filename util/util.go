package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadInput(filename string) (string, error) {
	inFile, err := os.Open(filename)
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
