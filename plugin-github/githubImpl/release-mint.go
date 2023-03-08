package githubImpl

import (
	"context"
	"fmt"

	"github.com/google/go-github/v49/github"

	zGithub "github.com/zostay/zedpm/pkg/github"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// ReleaseMintTask implements the /release/mint/github task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
	zGithub.Github
}

// CreateGithubPullRequest creates the PR on github for monitoring the test
// results for release testing. This will also be used to merge the release
// branch when testing passes.
func (s *ReleaseMintTask) CreateGithubPullRequest(ctx context.Context) error {
	owner, project, err := s.OwnerProject(ctx)
	if err != nil {
		return fmt.Errorf("failed getting owner/project information: %w", err)
	}

	branch, err := zGithub.ReleaseBranch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get release branch name: %w", err)
	}

	prName, err := zGithub.GetPropertyGithubReleaseName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get release name: %w", err)
	}

	body := fmt.Sprintf("Pull request to complete %q of project.", prName)
	if version := goals.GetPropertyReleaseVersion(ctx); version != "" {
		body = fmt.Sprintf("Pull request to complete release for v%s of project.", version)
	}

	_, _, err = s.Client().PullRequests.Create(ctx, owner, project, &github.NewPullRequest{
		Title: github.String(prName),
		Head:  github.String(branch),
		Base:  github.String(zGithub.TargetBranch(ctx)),
		Body:  github.String(body),
	})

	if err != nil {
		return fmt.Errorf("unable to create pull request: %w", err)
	}

	plugin.Logger(ctx,
		"operation", "CreateGithubPullRequest",
		"owner", owner,
		"project", project,
		"branch", branch,
		"pullRequestName", prName,
	).Info("Created Github pull request %q", prName)

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
