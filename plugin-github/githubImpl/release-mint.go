package githubImpl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v49/github"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/git"
	zGithub "github.com/zostay/zedpm/pkg/github"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/pkg/log"
	"github.com/zostay/zedpm/plugin"
)

// ReleaseMintTask implements the /release/mint/github task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
	zGithub.Github
}

// Setup configures the github and git clients.
func (s *ReleaseMintTask) Setup(ctx context.Context) error {
	return s.SetupGithubClient(ctx)
}

// CreateGithubPullRequest creates the PR on github for monitoring the test
// results for release testing. This will also be used to merge the release
// branch when testing passes.
func (s *ReleaseMintTask) CreateGithubPullRequest(ctx context.Context) error {
	logger := plugin.Logger(ctx, "operation", "CreateGithubPullRequest")
	logger.StartAction("CreateGithubPullRequest", "Creating Github pull request", "spin")

	owner, project, err := s.OwnerProject(ctx)
	if err != nil {
		logger.MarkAction("CreateGithubPullRequest", log.Fail)
		return format.WrapErr(err, "failed getting owner/project information")
	}

	logger.TickAction("CreateGithubPullRequest")

	branch, err := git.GetPropertyGitReleaseBranch(ctx)
	if err != nil {
		logger.MarkAction("CreateGithubPullRequest", log.Fail)
		return format.WrapErr(err, "failed to get release branch name")
	}

	logger.TickAction("CreateGithubPullRequest")

	prName, err := zGithub.GetPropertyGithubReleaseName(ctx)
	if err != nil {
		logger.MarkAction("CreateGithubPullRequest", log.Fail)
		return format.WrapErr(err, "failed to get release name")
	}

	logger.TickAction("CreateGithubPullRequest")

	body := fmt.Sprintf("Pull request to complete %q of project.", prName)
	if version, err := goals.GetPropertyReleaseVersion(ctx); err == nil {
		body = fmt.Sprintf("Pull request to complete release for v%s of project.", version)
	}

	logger.TickAction("CreateGithubPullRequest")

	targetBranch := git.GetPropertyGitTargetBranch(ctx)
	logger = logger.With(
		"owner", owner,
		"project", project,
		"branch", branch,
		"targetBranch", targetBranch,
		"pullRequestName", prName,
	)

	for retries := 3; retries > 0; retries-- {
		logger.MarkAction("CreateGithubPullRequest", log.Working)
		_, _, err = s.Client().PullRequests.Create(ctx, owner, project, &github.NewPullRequest{
			Title: github.String(prName),
			Head:  github.String(branch),
			Base:  github.String(targetBranch),
			Body:  github.String(body),
		})

		logger.TickAction("CreateGithubPullRequest")

		if err == nil {
			break
		}

		logger.MarkAction("CreateGithubPullRequest", log.Retry)
		// logger.Info("failed to create pull request; retrying in 5s", "error", err)

		<-time.After(5 * time.Second)
	}

	if err != nil {
		logger.MarkAction("CreateGithubPullRequest", log.Fail)
		return format.WrapErr(err, "unable to create pull request")
	}

	logger.MarkAction("CreateGithubPullRequest", log.Pass)

	return nil
}

// End configures CreateGithubPullRequest to run.
func (s *ReleaseMintTask) End(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  80,
			Action: plugin.OperationFunc(s.CreateGithubPullRequest),
		},
	}, nil
}
