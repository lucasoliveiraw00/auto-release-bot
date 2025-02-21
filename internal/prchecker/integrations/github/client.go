package github

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v69/github"
	"github.com/lucasoliveiraw00/auto-release-bot/config"
)

// NewClient inicializa e retorna um cliente autenticado do GitHub
func NewClient(ctx context.Context, cfg config.Config) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.GithubToken})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
