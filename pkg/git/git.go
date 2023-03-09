package git

import (
	"context"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/plugin"
)

// TODO IgnoreStatus should probably not be a thing. We can just use .gitignore.

// IgnoreStatus defines some global files to always ignore when checking for
// dirtiness.
var IgnoreStatus = map[string]struct{}{
	".session.vim": {},
}

// Git provides tools for working with a Git repository. It sets up client
// objects to work with the local repository, the remote repository, and the
// local work tree.
type Git struct {
	repo   *git.Repository
	remote *git.Remote
	wc     *git.Worktree
}

func ref(t, n string) plumbing.ReferenceName {
	return plumbing.ReferenceName(path.Join("refs", t, n))
}

func refSpec(r plumbing.ReferenceName) gitConfig.RefSpec {
	sr := string(r)
	return gitConfig.RefSpec(strings.Join([]string{sr, sr}, ":"))
}

func TargetBranch(ctx context.Context) string {
	if plugin.IsSet(ctx, "target_branch") {
		return plugin.GetString(ctx, "target_branch")
	}
	return "master"
}

func TargetBranchRefName(ctx context.Context) plumbing.ReferenceName {
	return ref("heads", TargetBranch(ctx))
}

func ReleaseBranchRefName(ctx context.Context) (plumbing.ReferenceName, error) {
	branch, err := GetPropertyGitReleaseBranch(ctx)
	if err != nil {
		return "", err
	}
	return ref("heads", branch), nil
}

func ReleaseBranchRefSpec(ctx context.Context) (gitConfig.RefSpec, error) {
	branchRefName, err := ReleaseBranchRefName(ctx)
	if err != nil {
		return "", err
	}
	return refSpec(branchRefName), nil
}

func ReleaseTagRefSpec(ctx context.Context) (gitConfig.RefSpec, error) {
	tag, err := GetPropertyGitReleaseTag(ctx)
	if err != nil {
		return "", err
	}
	tagRefName := ref("tags", tag)
	return refSpec(tagRefName), nil
}

func (g *Git) SetupGitRepo(ctx context.Context) error {
	l, err := git.PlainOpen(".")
	if err != nil {
		return format.WrapErr(err, "unable to open git repository at .")
	}

	g.repo = l

	r, err := g.repo.Remote("origin")
	if err != nil {
		return format.WrapErr(err, "unable to connect to remote origin")
	}

	g.remote = r

	w, err := g.repo.Worktree()
	if err != nil {
		return format.WrapErr(err, "unable to examine the working copy")
	}

	g.wc = w

	return nil
}

func (g *Git) Repository() *git.Repository {
	return g.repo
}

func (g *Git) Remote() *git.Remote {
	return g.remote
}

func (g *Git) Worktree() *git.Worktree {
	return g.wc
}
