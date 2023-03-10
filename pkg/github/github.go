package github

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"

	"github.com/zostay/zedpm/pkg/git"
)

// Github provides a client object for accessing the Github API.
type Github struct {
	git.Git
	gh *github.Client
}

func (g *Github) Client() *github.Client {
	return g.gh
}

func (g *Github) SetupGithubClient(ctx context.Context) error {
	err := g.Git.SetupGitRepo(ctx)
	if err != nil {
		return err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is missing")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	g.gh = github.NewClient(tc)

	return nil
}

var githubPrefixes = []string{
	"git@github.com:",
	"https://github.com/",
}

func (g *Github) OwnerProject(ctx context.Context) (string, string, error) {
	owner := GetPropertyGithubOwner(ctx)
	project := GetPropertyGithubProject(ctx)

	if owner != "" && project != "" {
		return owner, project, nil
	}

	if g.Remote() == nil {
		return owner, project, fmt.Errorf("unable to dtermine Github project and owner from git remote configuration: unable to load git remote client")
	}

	cfg := g.Remote().Config()
	if cfg == nil {
		return owner, project, fmt.Errorf("unable to determine Github project and owner from git remote configuration: no remote configuration found")
	}

	urls := cfg.URLs
	if len(urls) == 0 {
		return owner, project, fmt.Errorf("unable to determine Github project and owner from git remote configuration: no remote URLs found")
	}

	url := urls[0]
	for _, prefix := range githubPrefixes {
		if !strings.HasPrefix(url, prefix) {
			continue
		}

		urlPath := url[len(prefix):]
		if strings.HasSuffix(urlPath, ".git") {
			urlPath = urlPath[:len(urlPath)-len(".git")-1]
		}

		parts := strings.Split(urlPath, "/")
		if len(parts) == 2 {
			if owner == "" {
				owner = parts[0]
			}
			if project == "" {
				project = parts[1]
			}

			return owner, project, nil
		}
	}

	return owner, project, fmt.Errorf("unable to determing Github project and owner from git remote configuration: remote URL does not look like a github URL")
}
