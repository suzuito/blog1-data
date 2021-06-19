package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
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
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		if len(files) <= 0 {
			break
		}
		for _, file := range files {
			fmt.Println(*file.Filename, *file.Status)
		}
		fmt.Println(res)
		if res.NextPage == res.LastPage {
			break
		}
		page = res.NextPage
	}
}
