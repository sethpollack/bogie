package io

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"go.mozilla.org/sops/decrypt"
	"golang.org/x/oauth2"
)

type GitConfig struct {
	owner string
	repo  string
	path  string
}

type FileData struct {
	name string
	dir  bool
}

func (f *FileData) IsDir() bool {
	return f.dir
}

func (f *FileData) Name() string {
	return f.name
}

func ReadDir(dirname string) ([]FileData, error) {
	if ok := isValidUrl(dirname); !ok {
		return readDir(dirname)
	}

	return readGitDir(dirname)
}

func readDir(dirname string) ([]FileData, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	filedata := []FileData{}
	for _, f := range files {
		filedata = append(filedata, FileData{name: f.Name(), dir: f.IsDir()})
	}

	return filedata, nil
}

func getClient() (client *github.Client) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func readGitDir(dirname string) ([]FileData, error) {
	client := getClient()
	ctx := context.Background()
	conf := parseUrl(dirname)

	_, dc, _, err := client.Repositories.GetContents(ctx, conf.owner, conf.repo, conf.path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}

	files := []FileData{}
	for _, f := range dc {
		files = append(files, FileData{name: *f.Name, dir: *f.Type == "dir"})
	}

	return files, nil
}

func DecryptFile(filename, format string) ([]byte, error) {
	if ok := isValidUrl(filename); !ok {
		return decrypt.File(filename, format)
	}

	content, err := readGitInput(filename)
	if err != nil {
		return nil, err
	}

	return decrypt.Data([]byte(content), format)
}

func ReadInput(filename string) (string, error) {
	if ok := isValidUrl(filename); !ok {
		return readInput(filename)
	}

	return readGitInput(filename)
}

func readInput(filename string) (string, error) {
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

func readGitInput(path string) (string, error) {
	client := getClient()
	ctx := context.Background()
	conf := parseUrl(path)
	fc, _, _, err := client.Repositories.GetContents(ctx, conf.owner, conf.repo, conf.path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return "", err
	}

	return fc.GetContent()
}

func parseUrl(url string) GitConfig {
	split := strings.SplitN(url, "/", 6)
	return GitConfig{owner: split[3], repo: split[4], path: split[5]}
}

func isValidUrl(path string) bool {
	_, err := url.ParseRequestURI(path)
	return err == nil
}
