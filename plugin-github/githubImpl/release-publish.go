package githubImpl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v49/github"

	"github.com/zostay/zedpm/format"
	zGithub "github.com/zostay/zedpm/pkg/github"
	"github.com/zostay/zedpm/plugin"
)

// ReleasePublishTask implements the /release/publish/github task.
type ReleasePublishTask struct {
	plugin.TaskBoilerplate
	zGithub.Github
}

// CheckReadyForMerge ensures that all the required tests are passing.
func (f *ReleasePublishTask) CheckReadyForMerge(ctx context.Context) error {
	owner, project, err := f.OwnerProject(ctx)
	if err != nil {
		return format.WrapErr(err, "failed getting owner/project information")
	}

	branch, err := zGithub.ReleaseBranch(ctx)
	if err != nil {
		return format.WrapErr(err, "failed to get release branch name")
	}

	bp, _, err := f.Client().Repositories.GetBranchProtection(ctx, owner, project, zGithub.TargetBranch(ctx))
	if err != nil {
		return format.WrapErr(err, "unable to get branches %s", branch)
	}

	checks := bp.GetRequiredStatusChecks().Checks
	passage := make(map[string]bool, len(checks))
	for _, check := range checks {
		passage[check.Context] = false
	}

	crs, _, err := f.Client().Checks.ListCheckRunsForRef(ctx, owner, project, branch, &github.ListCheckRunsOptions{})
	if err != nil {
		return format.WrapErr(err, "unable to list check runs for branch %s", branch)
	}

	for _, run := range crs.CheckRuns {
		passage[run.GetName()] =
			run.GetStatus() == "completed" &&
				run.GetConclusion() == "success"
	}

	for k, v := range passage {
		if !v {
			return fmt.Errorf("cannot merge release branch because it has not passed check %q", k)
		}
	}

	plugin.Logger(ctx,
		"operation", "CheckReadyForMerge",
		"owner", owner,
		"project", project,
		"branch", branch,
	).Info("All Github required checks appear to be passing")

	return nil
}

// Check executes CheckReadyForMerge in a loop until either the Github checks
// pass or 15 minutes have elapsed, whichever comes first.
func (f *ReleasePublishTask) Check(ctx context.Context) error {
	const timeout = 15 * time.Minute
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var err error
	for {
		if ctx.Err() != nil {
			break
		}

		err = f.CheckReadyForMerge(ctx)
		if err != nil {
			<-time.After(30 * time.Second)
		}
	}

	return err
}

// MergePullRequest merges the PR into master.
func (f *ReleasePublishTask) MergePullRequest(ctx context.Context) error {
	owner, project, err := f.OwnerProject(ctx)
	if err != nil {
		return format.WrapErr(err, "failed getting owner/project information")
	}

	prs, _, err := f.Client().PullRequests.List(ctx, owner, project, &github.PullRequestListOptions{})
	if err != nil {
		return format.WrapErr(err, "unable to list pull requests")
	}

	branch, err := zGithub.ReleaseBranch(ctx)
	if err != nil {
		return format.WrapErr(err, "failed to get release branch name")
	}

	prId := 0
	for _, pr := range prs {
		if pr.Head.GetRef() == branch {
			prId = pr.GetNumber()
			break
		}
	}

	if prId == 0 {
		return fmt.Errorf("cannot find pull request for branch %s", branch)
	}

	m, _, err := f.Client().PullRequests.Merge(ctx, owner, project, prId, "Merging release branch.", &github.PullRequestOptions{})
	if err != nil {
		return format.WrapErr(err, "unable to merge pull request %d", prId)
	}

	if !m.GetMerged() {
		return fmt.Errorf("failed to merge pull request %d", prId)
	}

	plugin.Logger(ctx,
		"operation", "MergePullRequest",
		"owner", owner,
		"project", project,
		"branch", branch,
		"pullRequestID", prId,
	).Info("Merged the pull request into the target branch.")

	return nil
}

// CreateRelease creates a release on github for the release.
func (f *ReleasePublishTask) CreateRelease(ctx context.Context) error {
	owner, project, err := f.OwnerProject(ctx)
	if err != nil {
		return format.WrapErr(err, "failed getting owner/project information")
	}

	tag, err := zGithub.ReleaseTag(ctx)
	if err != nil {
		return format.WrapErr(err, "failed to get release tag name")
	}

	releaseName, err := zGithub.GetPropertyGithubReleaseName(ctx)
	if err != nil {
		return format.WrapErr(err, "failed to get release name")
	}

	changesInfo := zGithub.ReleaseDescription(ctx)
	_, _, err = f.Client().Repositories.CreateRelease(ctx, owner, project,
		&github.RepositoryRelease{
			TagName:              github.String(tag),
			Name:                 github.String(releaseName),
			Body:                 github.String(changesInfo),
			Draft:                github.Bool(false),
			Prerelease:           github.Bool(false),
			GenerateReleaseNotes: github.Bool(false),
			MakeLatest:           github.String("true"),
		},
	)

	if err != nil {
		return format.WrapErr(err, "failed to create release %q", releaseName)
	}

	plugin.Logger(ctx,
		"owner", owner,
		"project", project,
		"tag", tag,
		"releaseName", releaseName,
	).Info("Created a release named %q", releaseName)

	return nil
}

// Run configures MergePullRequest and CreateRelease to run.
func (f *ReleasePublishTask) Run(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  70,
			Action: plugin.OperationFunc(f.MergePullRequest),
		},
		{
			Order:  75,
			Action: plugin.OperationFunc(f.CreateRelease),
		},
	}, nil
}
