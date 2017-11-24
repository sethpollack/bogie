package io

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"go.mozilla.org/sops/decrypt"
	"golang.org/x/oauth2"
)

type repoInfo struct {
	owner string
	repo  string
	path  string
}

type fileInfo struct {
	name string
	dir  bool
	size int64
	time time.Time
}

func (f *fileInfo) IsDir() bool        { return f.dir }
func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) Size() int64        { return f.size }
func (f *fileInfo) ModTime() time.Time { return f.time }
func (f *fileInfo) Sys() interface{}   { return nil }
func (f *fileInfo) Mode() os.FileMode {
	if f.IsDir() {
		return os.ModePerm | os.ModeDir
	}
	return os.FileMode(0644)
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	if ok := isValidUrl(dirname); !ok {
		return ioutil.ReadDir(dirname)
	}

	return readRepoDir(dirname)
}

func ReadFile(filename string) ([]byte, error) {
	if ok := isValidUrl(filename); !ok {
		return readFile(filename)
	}

	return readRepoFile(filename)
}

func DecryptFile(filename, format string) ([]byte, error) {
	if ok := isValidUrl(filename); !ok {
		return decrypt.File(filename, format)
	}

	content, err := readRepoFile(filename)
	if err != nil {
		return nil, err
	}

	return decrypt.Data([]byte(content), format)
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

func readRepoDir(dirname string) ([]os.FileInfo, error) {
	client := getClient()
	ctx := context.Background()
	conf := parseUrl(dirname)

	_, dc, _, err := client.Repositories.GetContents(ctx, conf.owner, conf.repo, conf.path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}

	files := []os.FileInfo{}
	for _, f := range dc {
		files = append(files, &fileInfo{name: *f.Name, size: int64(*f.Size), dir: *f.Type == "dir"})
	}

	return files, nil
}

func readFile(filename string) ([]byte, error) {
	inFile, err := os.Open(filename)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to open %s\n%v", filename, err)
	}
	defer inFile.Close()

	bytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		err = fmt.Errorf("read failed for %s\n%v", filename, err)
		return []byte{}, err
	}

	return bytes, nil
}

func readRepoFile(path string) ([]byte, error) {
	client := getClient()
	conf := parseUrl(path)

	fc, _, _, err := client.Repositories.GetContents(
		context.Background(),
		conf.owner,
		conf.repo,
		conf.path,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil {
		return []byte{}, err
	}

	content, err := fc.GetContent()

	return []byte(content), err
}

func parseUrl(url string) repoInfo {
	split := strings.SplitN(url, "/", 6)
	return repoInfo{owner: split[3], repo: split[4], path: split[5]}
}

func isValidUrl(path string) bool {
	_, err := url.ParseRequestURI(path)
	return err == nil
}
