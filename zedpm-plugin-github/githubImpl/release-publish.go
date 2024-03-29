package githubImpl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v49/github"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/git"
	zGithub "github.com/zostay/zedpm/pkg/github"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/pkg/log"
	"github.com/zostay/zedpm/plugin"
)

// ReleasePublishTask implements the /release/publish/github task.
type ReleasePublishTask struct {
	plugin.TaskBoilerplate
	zGithub.Github
}

// Setup configures the github and git clients.
func (s *ReleasePublishTask) Setup(ctx context.Context) error {
	return s.SetupGithubClient(ctx)
}

// CheckReadyForMerge ensures that all the required tests are passing.
func (f *ReleasePublishTask) CheckReadyForMerge(ctx context.Context) error {
	logger := plugin.Logger(ctx, "operation", "CheckReadyForMerge")
	logger.StartAction("CheckReadyForMerge", "Checking if pull request is ready to merge", "spin")

	owner, project, err := f.OwnerProject(ctx)
	if err != nil {
		logger.MarkAction("CheckReadyForMerge", log.Fail)
		return format.WrapErr(err, "failed getting owner/project information")
	}

	logger = logger.With("owner", owner, "project", project)
	logger.TickAction("CheckReadyForMerge")

	branch, err := git.GetPropertyGitReleaseBranch(ctx)
	if err != nil {
		logger.MarkAction("CheckReadyForMerge", log.Fail)
		return format.WrapErr(err, "failed to get release branch name")
	}

	logger = logger.With("branch", branch)
	logger.TickAction("CheckReadyForMerge")

	bp, _, err := f.Client().Repositories.GetBranchProtection(ctx, owner, project, git.GetPropertyGitTargetBranch(ctx))
	if err != nil {
		if strings.Contains(err.Error(), "branch is not protected") {
			logger.MarkAction("CheckReadyForMerge", log.Pass)
			logger.Info("Branch is ready to merge: no protection")
			return nil
		}
		logger.MarkAction("CheckReadyForMerge", log.Fail)
		return format.WrapErr(err, "unable to get branches %s", branch)
	}

	checks := bp.GetRequiredStatusChecks().Checks
	passage := make(map[string]bool, len(checks))
	for _, check := range checks {
		passage[check.Context] = false
	}

	crs, _, err := f.Client().Checks.ListCheckRunsForRef(ctx, owner, project, branch, &github.ListCheckRunsOptions{})
	if err != nil {
		logger.MarkAction("CheckReadyForMerge", log.Fail)
		return format.WrapErr(err, "unable to list check runs for branch %s", branch)
	}

	for _, run := range crs.CheckRuns {
		passage[run.GetName()] =
			run.GetStatus() == "completed" &&
				run.GetConclusion() == "success"
	}

	for k, v := range passage {
		if !v {
			logger.MarkAction("CheckReadyForMerge", log.Fail)
			return fmt.Errorf("cannot merge release branch because it has not passed check %q", k)
		}
	}

	logger.MarkAction("CheckReadyForMerge", log.Pass)
	logger.Info("Branch is ready for merge: All Github required checks appear to be passing")

	return nil
}

// Check executes CheckReadyForMerge in a loop until either the Github checks
// pass or 15 minutes have elapsed, whichever comes first.
func (f *ReleasePublishTask) Check(ctx context.Context) error {
	// TODO Make this timeout into a property like github.publishWaitTimeout
	const timeout = 1 * time.Minute
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Restart the check every few seconds, but not if it's still running for
	// some reason.
	var waiting, running bool
	errCh := make(chan error)
	startReadinessCheck := func() {
		if !waiting && !running {
			running, waiting = true, true
			go func() {
				errCh <- f.CheckReadyForMerge(ctx)
			}()
		}
	}

	startReadinessCheck()

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			if lastErr != nil {
				return lastErr
			}
			return ctx.Err()

		case <-time.After(30 * time.Second):
			// TODO Make the retry delay here configurable by property
			waiting = false
			startReadinessCheck()
			// TODO If running == true for too many of these delays, maybe we want to cancel early?
			continue

		case lastErr = <-errCh:
			running = false
			startReadinessCheck()
			if lastErr == nil {
				return nil
			}
		}
	}
}

// MergePullRequest merges the PR into master.
func (f *ReleasePublishTask) MergePullRequest(ctx context.Context) error {
	logger := plugin.Logger(ctx, "operation", "MergePullRequest")
	logger.StartAction("MergePullRequest", "Merging pull request", "spin")

	owner, project, err := f.OwnerProject(ctx)
	if err != nil {
		logger.MarkAction("MergePullRequest", log.Fail)
		return format.WrapErr(err, "failed getting owner/project information")
	}

	logger = logger.With("owner", owner, "project", project)
	logger.TickAction("MergePullRequest")

	prs, _, err := f.Client().PullRequests.List(ctx, owner, project, &github.PullRequestListOptions{})
	if err != nil {
		logger.MarkAction("MergePullRequest", log.Fail)
		return format.WrapErr(err, "unable to list pull requests")
	}

	logger.TickAction("MergePullRequest")

	branch, err := git.GetPropertyGitReleaseBranch(ctx)
	if err != nil {
		logger.MarkAction("MergePullRequest", log.Fail)
		return format.WrapErr(err, "failed to get release branch name")
	}

	logger = logger.With("branch", branch)
	logger.TickAction("MergePullRequest")

	prId := 0
	for _, pr := range prs {
		if pr.Head.GetRef() == branch {
			prId = pr.GetNumber()
			break
		}
	}

	if prId == 0 {
		logger.MarkAction("MergePullRequest", log.Fail)
		return fmt.Errorf("cannot find pull request for branch %s", branch)
	}

	logger = logger.With("pullRequestID", prId)
	logger.TickAction("MergePullRequest")

	m, _, err := f.Client().PullRequests.Merge(ctx, owner, project, prId, "Merging release branch.", &github.PullRequestOptions{})
	if err != nil {
		logger.MarkAction("MergePullRequest", log.Fail)
		return format.WrapErr(err, "unable to merge pull request %d", prId)
	}

	if !m.GetMerged() {
		logger.MarkAction("MergePullRequest", log.Fail)
		return fmt.Errorf("failed to merge pull request %d", prId)
	}

	logger.MarkAction("MergePullRequest", log.Pass)
	logger.Info("Merged the pull request into the target branch.")

	return nil
}

// CreateRelease creates a release on github for the release.
func (f *ReleasePublishTask) CreateRelease(ctx context.Context) error {
	logger := plugin.Logger(ctx, "operation", "CreateRelease")
	logger.StartAction("CreateRelease", "Creating a Github release", "spin")

	owner, project, err := f.OwnerProject(ctx)
	if err != nil {
		logger.MarkAction("CreateRelease", log.Fail)
		return format.WrapErr(err, "failed getting owner/project information")
	}

	logger = logger.With("owner", owner, "project", project)
	logger.TickAction("CreateRelease")

	tag, err := git.GetPropertyGitReleaseTag(ctx)
	if err != nil {
		logger.MarkAction("CreateRelease", log.Fail)
		return format.WrapErr(err, "failed to get release tag name")
	}

	logger = logger.With("tag", tag)
	logger.TickAction("CreateRelease")

	releaseName, err := zGithub.GetPropertyGithubReleaseName(ctx)
	if err != nil {
		logger.MarkAction("CreateRelease", log.Fail)
		return format.WrapErr(err, "failed to get release name")
	}

	logger = logger.With("releaseName", releaseName)
	logger.TickAction("CreateRelease")

	changesInfo := goals.GetPropertyReleaseDescription(ctx)
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
		logger.MarkAction("CreateRelease", log.Fail)
		return format.WrapErr(err, "failed to create release %q", releaseName)
	}

	logger.MarkAction("CreateRelease", log.Pass)

	return nil
}

// Run configures MergePullRequest and CreateRelease to run.
func (f *ReleasePublishTask) End(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  60,
			Action: plugin.OperationFunc(f.MergePullRequest),
		},
		{
			Order:  80,
			Action: plugin.OperationFunc(f.CreateRelease),
		},
	}, nil
}
