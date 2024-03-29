package gitImpl

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/zostay/zedpm/format"
	zGit "github.com/zostay/zedpm/pkg/git"
	"github.com/zostay/zedpm/plugin"
)

// ReleaseMintTask implements the /release/mint/git task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
	zGit.Git
}

// Setup initializes the git client and related objects.
func (s *ReleaseMintTask) Setup(ctx context.Context) error {
	return s.SetupGitRepo(ctx)
}

// IsDirty returns true if we consider the tree dirty. We do not consider
// Untracked to dirty the directory and we also ignore some filenames that are
// in the global .gitignore and not in the local .gitignore.
func IsDirty(status git.Status) bool {
	for fn, fstat := range status {
		if _, ignorable := zGit.IgnoreStatus[fn]; ignorable {
			continue
		}

		if fstat.Worktree != git.Unmodified && fstat.Worktree != git.Untracked {
			return true
		}

		if fstat.Staging != git.Unmodified && fstat.Staging != git.Untracked {
			return true
		}
	}
	return false
}

// CheckGitCleanliness ensures that the current git repository is clean and that
// we are on the correct branch from which to trigger a release.
func (s *ReleaseMintTask) CheckGitCleanliness(ctx context.Context) error {
	logger := plugin.Logger(ctx,
		"operation", "CheckGitCleanliness",
		"task", "/release/mint/git",
	)
	logger.Info("Finding the HEAD reference")

	headRef, err := s.Repository().Head()
	if err != nil {
		return format.WrapErr(err, "unable to find HEAD")
	}

	if headRef.Name() != zGit.TargetBranchRefName(ctx) {
		return fmt.Errorf("you must checkout %s to release", zGit.TargetBranch(ctx))
	}

	logger = logger.With("headRef", headRef.String())
	logger.Info("Finding the remote master reference")

	remoteRefs, err := s.Remote().List(&git.ListOptions{})
	if err != nil {
		return format.WrapErr(err, "unable to list remote git references")
	}

	var masterRef *plumbing.Reference
	for _, ref := range remoteRefs {
		if ref.Name() == zGit.TargetBranchRefName(ctx) {
			masterRef = ref
			break
		}
	}

	logger = logger.With("masterRef", masterRef.String())
	logger.Info("Checking if local master reference matches remote")

	if headRef.Hash() != masterRef.Hash() {
		return fmt.Errorf("local copy differs from remote, you need to push or pull")
	}

	logger.Info("Checking that the local copy is clean")

	stat, err := s.Worktree().Status()
	if err != nil {
		return format.WrapErr(err, "unable to check working copy status")
	}

	ignoreDirty := GetPropertyGitIgnoreDirty(ctx)
	if !ignoreDirty && IsDirty(stat) {
		return fmt.Errorf("your working copy is dirty")
	}

	logger.Info("Git working tree is clean for release")

	return nil
}

// Check calls CheckGitCleanliness.
func (s *ReleaseMintTask) Check(ctx context.Context) error {
	return s.CheckGitCleanliness(ctx)
}

// MakeReleaseBranch creates the branch that will be used to manage the release.
func (s *ReleaseMintTask) MakeReleaseBranch(ctx context.Context) error {
	headRef, err := s.Repository().Head()
	if err != nil {
		return format.WrapErr(err, "unable to retrieve the HEAD ref")
	}

	branchRefName, err := zGit.ReleaseBranchRefName(ctx)
	if err != nil {
		return format.WrapErr(err, "unable to determine release branch references")
	}

	branch, _ := zGit.GetPropertyGitReleaseBranch(ctx)
	err = s.Worktree().Checkout(&git.CheckoutOptions{
		Hash:   headRef.Hash(),
		Branch: branchRefName,
		Create: true,
	})
	if err != nil {
		return format.WrapErr(err, "unable to checkout branch %s", branch)
	}

	plugin.ForCleanup(ctx, func() {
		_ = s.Repository().Storer.RemoveReference(branchRefName)
	})
	plugin.ForCleanup(ctx, func() {
		_ = s.Worktree().Checkout(&git.CheckoutOptions{
			Branch: zGit.TargetBranchRefName(ctx),
		})
	})

	plugin.Logger(ctx,
		"operation", "MakeReleaseBranch",
		"headRef", headRef,
		"branchRefName", branchRefName,
		"branch", branch,
	).Info("Created git branch for managing the release")

	return nil
}

// Run sets up the MakeReleaseBranch operation.
func (s *ReleaseMintTask) Run(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  30,
			Action: plugin.OperationFunc(s.MakeReleaseBranch),
		},
	}, nil
}

// AddAndCommit adds changes made as part of the release process to the release
// branch.
func (s *ReleaseMintTask) AddAndCommit(ctx context.Context) error {
	logger := plugin.Logger(ctx)
	addedFiles := plugin.ListAdded(ctx)
	for _, fn := range addedFiles {
		_, err := s.Worktree().Add(fn)
		if err != nil {
			return format.WrapErr(err, "error adding file %s to git", fn)
		}

		logger.Info("Adding file to git", "filename", fn)
	}

	version := plugin.GetString(ctx, "release.version")
	msg := "releng: v" + version
	_, err := s.Worktree().Commit("releng: v"+version, &git.CommitOptions{})
	if err != nil {
		return format.WrapErr(err, "error committing changes to git")
	}

	plugin.Logger(ctx,
		"count", len(addedFiles),
		"version", version,
		"message", msg,
	).Info("Added files and committing changes to git")

	return nil
}

// PushReleaseBranch pushes the release branch to github for release testing.
func (s *ReleaseMintTask) PushReleaseBranch(ctx context.Context) error {
	branchRefSpec, err := zGit.ReleaseBranchRefSpec(ctx)
	if err != nil {
		return format.WrapErr(err, "unable to determine the ref spec")
	}

	err = s.Repository().Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{branchRefSpec},
	})
	if err != nil {
		return format.WrapErr(err, "error pushing changes to github branch %q", branchRefSpec.String())
	}

	plugin.ForCleanup(ctx, func() {
		_ = s.Remote().Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs:   []config.RefSpec{branchRefSpec},
			Prune:      true,
		})
	})

	plugin.Logger(ctx,
		"branchRefSpec", branchRefSpec,
	).Info("Pushed release branch to remote repository")

	return nil
}

// End sets up the AddAndCommit and PushReleaseBranch operations.
func (s *ReleaseMintTask) End(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  70,
			Action: plugin.OperationFunc(s.AddAndCommit),
		},
		{
			Order:  75,
			Action: plugin.OperationFunc(s.PushReleaseBranch),
		},
	}, nil
}
