package main

import (
	"fmt"
	"repowatch/repowatch"
)

func main() {
	config, err := repowatch.ParseConfig()
	if err != nil {
		panic(err)
	}

	pullRequests := repowatch.FetchPullRequests(config)

	for _, pr := range pullRequests {
		fmt.Println(pr.String())
	}
}
