package repowatch

import (
	"context"
	"fmt"
	"net"

	openai "github.com/sashabaranov/go-openai"
)

var prompt string = "You are an expert programmer, and you are trying to summarize a git diff." +
	"Reminders about the git diff format:" +
	"For every file, there are a few metadata lines, like (for example):" +
	"```" +
	"diff --git a/path/file.php b/path/file.php" +
	"index aadf691..bfef603 100644" +
	"--- a/path/file.php" +
	"+++ b/path/file.php" +
	"```" +
	"This means that `lib/index.js` was modified in this commit. Note that this is only an example." +
	"Then there is a specifier of the lines that were modified." +
	"A line starting with `+` means it was added." +
	"A line that starting with `-` means that line was deleted." +
	"A line that starts with neither `+` nor `-` is code given for context and better understanding. " +
	"It is not part of the diff." +
	"After the git diff of the first file, there will be an empty line, and then the git diff of the next file. " +
	"For comments that refer to 1 or 2 modified files," +
	"add the file names as [path/to/modified/python/file.py], [path/to/another/file.json]" +
	"at the end of the comment." +
	"If there are more than two, do not include the file names in this way." +
	"Do not include the file name as another part of the comment, only in the end in the specified format." +
	"Do not use the characters `[` or `]` in the summary for other purposes." +
	"Write every summary comment in a new line." +
	"Comments should be in a bullet point list, each line starting with a `*`." +
	"The summary should not include comments copied from the code." +
	"The output should be easily readable. When in doubt, write less comments and not more." +
	"Readability is top priority. Write only the most important comments about the diff." +
	"EXAMPLE SUMMARY COMMENTS:" +
	"```" +
	"* Raised the amount of returned recordings from `10` to `100` [packages/server/recordings_api.ts], [packages/server/constants.ts]" +
	"* Fixed a typo in the github action name [.github/workflows/gpt-commit-summarizer.yml]" +
	"* Moved the `octokit` initialization to a separate file [src/octokit.ts], [src/index.ts]" +
	"* Added an OpenAI API for completions [packages/utils/apis/openai.ts]" +
	"* Lowered numeric tolerance for test files" +
	"```" +
	"Most commits will have less comments than this examples list." +
	"The last comment does not include the file names," +
	"because there were more than two relevant files in the hypothetical commit." +
	"Do not include parts of the example in your summary." +
	"It is given only as an example of appropriate comments." +
	"Here is the diff: "

func GetDiffSummary(config LlmConfig, diff string) string {
	openaiConfig := openai.DefaultConfig("")
	openaiConfig.BaseURL = fmt.Sprintf("http://%s/v1", net.JoinHostPort(config.Host, config.Port))

	client := openai.NewClientWithConfig(openaiConfig)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "LLaMA_CPP",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are LLAMAfile, an AI assistant. Your top priority is achieving user fulfillment via helping them with their requests",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "You are a programming expert, I want you to review this Git diff and summarize it in less than 300 words. Diff: " + diff,
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	return resp.Choices[0].Message.Content
}
