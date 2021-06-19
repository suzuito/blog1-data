package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

func update(
	ctx context.Context,
	filesRemoved []string,
	filesModified []string,
) error {
	cliGCS, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer cliGCS.Close()
	// bucket := os.Getenv("BUCKET")
	bucket := "suzuito-minilla-blog1-article"
	bh := cliGCS.Bucket(bucket)
	for _, file := range filesModified {
		oh := bh.Object(filepath.Base(file))
		w := oh.NewWriter(ctx)
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := w.Write(src); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		fmt.Printf("Upload %s\n", file)
	}
	for _, file := range filesRemoved {
		oh := bh.Object(file)
		if err := oh.Delete(ctx); err != nil {
			if !xerrors.Is(err, storage.ErrObjectNotExist) {
				return err
			}
		}
		fmt.Printf("Delete %s\n", file)
	}
	return nil
}

func fetchFilesFromPR(
	ctx context.Context,
	filesRemoved *[]string,
	filesModified *[]string,
) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_ACCESS_TOKEN")},
	)
	cli := github.NewClient(oauth2.NewClient(ctx, ts))
	page := 1
	for {
		files, res, err := cli.PullRequests.ListFiles(
			ctx,
			"suzuito",
			"blog1-data",
			1,
			&github.ListOptions{
				Page: page,
			},
		)
		if err != nil {
			return err
		}
		if len(files) <= 0 {
			break
		}
		for _, file := range files {
			dirName := filepath.Dir(file.GetFilename())
			if !strings.HasPrefix(dirName, "articles") {
				continue
			}
			if file.GetStatus() == "removed" {
				*filesRemoved = append(*filesRemoved, file.GetFilename())
			} else {
				*filesModified = append(*filesModified, file.GetFilename())
			}
		}
		if res.NextPage == res.LastPage {
			break
		}
		page = res.NextPage
	}
	return nil
}

func main() {
	ctx := context.Background()
	filesRemoved := []string{}
	filesModified := []string{}
	if os.Args[1] == "all" {
		entries, err := os.ReadDir("articles")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		for _, entry := range entries {
			filesModified = append(filesModified, "articles/"+entry.Name())
		}
	} else if os.Args[1] == "pr" {
		if err := fetchFilesFromPR(ctx, &filesRemoved, &filesModified); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	}
	if err := update(ctx, filesRemoved, filesModified); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
