package main

import (
	"repowatch/repowatch"
)

func main() {
	config, err := repowatch.ParseConfig()
	if err != nil {
		panic(err)
	}

	pullRequests := repowatch.FetchPullRequests(config)
	repowatch.TablePrint(pullRequests)
}
