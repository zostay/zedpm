package gitImpl

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"

	"github.com/zostay/zedpm/plugin"
	zGit "github.com/zostay/zedpm/plugin-git/pkg/git"
)

// ReleasePublishTask implements the /release/publish/git task.
type ReleasePublishTask struct {
	plugin.TaskBoilerplate
	zGit.Git
}

// Setup initializes the git client and related objects.
func (f *ReleasePublishTask) Setup(ctx context.Context) error {
	return f.SetupGitRepo(ctx)
}

// TagRelease creates and pushes a tag for the newly merged release on master.
func (f *ReleasePublishTask) TagRelease(ctx context.Context) error {
	err := f.Worktree().Checkout(&git.CheckoutOptions{
		Branch: zGit.TargetBranchRefName(ctx),
	})
	if err != nil {
		return fmt.Errorf("unable to switch to %s branch: %w", zGit.TargetBranch(ctx), err)
	}

	headRef, err := f.Repository().Head()
	if err != nil {
		return fmt.Errorf("unable to get HEAD ref of %s branch: %w", zGit.TargetBranch(ctx), err)
	}

	tag, err := zGit.GetPropertyGitReleaseTag(ctx)
	if err != nil {
		return fmt.Errorf("unable to determine release tag: %w", err)
	}

	head := headRef.Hash()
	_, err = f.Repository().CreateTag(tag, head, &git.CreateTagOptions{
		Message: fmt.Sprintf("Release tag %q", tag),
	})
	if err != nil {
		return fmt.Errorf("unable to tag release %s: %w", tag, err)
	}

	plugin.ForCleanup(ctx, func() { _ = f.Repository().DeleteTag(tag) })

	tagRefSpec, err := zGit.ReleaseTagRefSpec(ctx)
	if err != nil {
		return fmt.Errorf("unable to determine release tag ref spec: %w", err)
	}

	err = f.Repository().Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{tagRefSpec},
	})
	if err != nil {
		return fmt.Errorf("unable to push tags to origin: %w", err)
	}

	plugin.ForCleanup(ctx, func() {
		_ = f.Remote().Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs:   []config.RefSpec{tagRefSpec},
			Prune:      true,
		})
	})

	plugin.Logger(ctx,
		"headRef", headRef,
		"tag", tag,
		"head", head,
		"tagRefSpec", tagRefSpec,
	).Info("Creating release tag %q and pushing to remote.", tag)

	return nil
}

// End sets up the TagRelease operation to run.
func (f *ReleasePublishTask) End(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  75,
			Action: plugin.OperationFunc(f.TagRelease),
		},
	}, nil
}
