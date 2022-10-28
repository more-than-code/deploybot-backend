package util

import (
	"archive/tar"
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func TarFiles(dir string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	fileSystem := os.DirFS(dir)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.IsDir() {
			return nil
		}

		bytes, err := os.ReadFile(dir + path)
		if err != nil {
			log.Fatal(err)
		}

		hdr := &tar.Header{
			Name: path,
			Mode: 0600,
			Size: int64(len(bytes)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write(bytes); err != nil {
			log.Fatal(err)
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	return &buf, nil
}

func CloneRepo(repoName, cloneUrl, username, token string) error {
	err := os.RemoveAll(repoName)

	if err != nil {
		return err
	}

	log.Println(cloneUrl, username, token)
	_, err = git.PlainClone(repoName, false, &git.CloneOptions{
		URL:               cloneUrl,
		Progress:          os.Stdout,
		RecurseSubmodules: 1,
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})

	if err != nil {
		return err
	}

	return nil
}