package repowatch

import (
  "strings"
  "time"

  "github.com/fatih/color"
  "github.com/rodaine/table"
)

func TablePrint(config *Config, prs []*PullRequest) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Repo", "State", "Title", "By", "Reviewers", "Created", "Link")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, pr := range prs {
		reviewers := strings.Join(pr.RequestedReviewers, "\n")
		state := pr.State


		if (pr.Draft) {
			state = state + " (Draft)"
		}

		tbl.AddRow(pr.RepoName, state, pr.Title, pr.CreatedBy, reviewers, pr.CreatedAt.Format(time.Stamp))
	}

	tbl.Print()
}
