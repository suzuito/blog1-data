package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/github"
	"github.com/suzuito/blog1-go/entity/model"
	"github.com/suzuito/blog1-go/inject"
	"github.com/suzuito/blog1-go/setting"
	"github.com/suzuito/blog1-go/usecase"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

func update(
	ctx context.Context,
	u usecase.Usecase,
	filesRemoved []string,
	filesModified []string,
) error {
	cliGCS, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer cliGCS.Close()
	bucket := os.Getenv("BUCKET")
	bh := cliGCS.Bucket(bucket)
	for _, file := range filesModified {
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		article := model.Article{}
		if err := u.ConvertMD(ctx, src, &article, &[]byte{}); err != nil {
			return err
		}
		oh := bh.Object(fmt.Sprintf("%s.md", article.ID))
		w := oh.NewWriter(ctx)
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
	prNumber int,
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
			prNumber,
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
	env, err := setting.NewEnvironment()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	gdeps, gcloseFunc, err := inject.NewGlobalDepends(ctx, env)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	defer gcloseFunc()
	cdeps, ccloseFunc, err := inject.NewContextDepends(ctx, env)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	defer ccloseFunc()
	u := usecase.NewImpl(env, cdeps.DB, cdeps.Storage, gdeps.MDConverter)
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
		prNumberString := os.Getenv("PR_NUMBER")
		prNumber, err := strconv.Atoi(prNumberString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		if err := fetchFilesFromPR(ctx, prNumber, &filesRemoved, &filesModified); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	}
	if err := update(ctx, u, filesRemoved, filesModified); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
