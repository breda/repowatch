package repowatch

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
)

type PullRequest struct {
	Draft              bool
	RepoName           string
	State              string // open,closed,merged
	Title              string
	CreatedBy          string
	CreatedAt          time.Time
	RequestedReviewers []string // List of teams requested for review
	Link			   string
	Summary            string
}

func FetchPullRequests(config *Config) []*PullRequest {
	ret := make([]*PullRequest, 0)

	client := github.NewClient(nil).WithAuthToken(config.Github.Token)

	for _, repoDef := range config.Repos {
		nameParts := strings.Split(repoDef.Name, "/")
		owner, repo := nameParts[0], nameParts[1]

		// Fetch pull requests
		githubPRs, resp, err := client.PullRequests.List(
			context.Background(),
			owner,
			repo,
			&github.PullRequestListOptions{
				State: "all",
			},
		)

		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			panic("github api did not return a 200 http response")
		}

		// Convert to our format
		for _, prData := range githubPRs {
			if len(ret) == repoDef.LimitNum {
				break
			}

			if prOlderThanLimit(prData.CreatedAt.Time, repoDef.LimitDays) {
				continue;
			}

			pr := &PullRequest{
				Draft:              *prData.Draft,
				RepoName:           repoDef.Name,
				State:              *prData.State,
				Title:              *prData.Title,
				CreatedBy:          *prData.User.Login,
				CreatedAt:          prData.CreatedAt.Time,
				RequestedReviewers: make([]string, 0),
			}

			if (config.Features.Summary) {
				diff := getPullRequestDiff(client, owner, repo, *prData.Number)
				pr.Summary = GetDiffSummary(config.LlmConfig, diff)
			}

			for _, requestedTeam := range prData.RequestedTeams {
				pr.RequestedReviewers = append(pr.RequestedReviewers, *requestedTeam.Name)
			}

			ret = append(ret, pr)
		}
	}

	return ret
}

func prOlderThanLimit(prCreationDate time.Time, limitDays int) bool {
	now := time.Now()

	duration := time.Duration(limitDays) * 24 * time.Hour
	threshold := now.Add(-duration)

	return prCreationDate.Before(threshold)
}

func getPullRequestDiff(client *github.Client, owner string, repo string, number int) string {
	diff, resp, err := client.PullRequests.GetRaw(context.Background(), owner, repo, number, github.RawOptions{Type: github.Diff})

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic("github api did not return a 200 http response")
	}

	return diff
}

func (pr *PullRequest) String() string {
	var format string = `Repository Name: %s
State: %s (draft: %s)
Title: %s
Created By: %s at %s
Requested Reviewers: %s
Summary: %s
`
	return fmt.Sprintf(
		format,
		pr.RepoName,
		pr.State,
		strconv.FormatBool(pr.Draft),
		pr.Title,
		pr.CreatedBy,
		pr.CreatedAt.Format("2006-1-2 10:00"),
		strings.Join(pr.RequestedReviewers, ","),
		pr.Summary,
	)
}
