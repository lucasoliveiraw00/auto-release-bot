package github

import (
	"context"
	"log"

	"github.com/google/go-github/v69/github"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
)

// Busca status de checks do PR
func CheckPRStatus(ctx context.Context, cfg config.Config, client *github.Client, ref string) (int, int) {
	status := "completed"
	opts := &github.ListCheckRunsOptions{Status: &status}

	checkRuns, _, err := client.Checks.ListCheckRunsForRef(ctx, cfg.GithubOwner, cfg.GithubRepo, ref, opts)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar checks para %s: %v", ref, err)
		return 0, 0
	}

	successfulChecks := 0
	failedChecks := 0

	for _, check := range checkRuns.CheckRuns {
		switch check.GetConclusion() {
		case "success":
			successfulChecks++
		case "failure":
			failedChecks++
		}
	}

	return successfulChecks, failedChecks
}

// Busca status de aprovação e revisão dos PRs
func GetApprovalStatus(ctx context.Context, cfg config.Config, client *github.Client, prNumber int) (bool, bool, int) {
	hasPendingChanges, hasUnresolvedComments := false, false
	approvalCount := 0

	reviews, _, _ := client.PullRequests.ListReviews(ctx, cfg.GithubOwner, cfg.GithubRepo, prNumber, nil)
	lastReviewByUser := make(map[string]string)

	for _, review := range reviews {
		lastReviewByUser[review.GetUser().GetLogin()] = review.GetState()
	}

	for _, state := range lastReviewByUser {
		switch state {
		case "CHANGES_REQUESTED":
			hasPendingChanges = true
		case "COMMENTED":
			hasUnresolvedComments = true
		case "APPROVED":
			approvalCount++
		}
	}

	return hasPendingChanges, hasUnresolvedComments, approvalCount
}

// Verifica se o PR tem conflitos
func HasMergeConflicts(ctx context.Context, cfg config.Config, client *github.Client, prNumber int) bool {
	pr, _, err := client.PullRequests.Get(ctx, cfg.GithubOwner, cfg.GithubRepo, prNumber)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar detalhes do PR #%d: %v", prNumber, err)
		return false
	}

	// O campo Mergeable pode ser nil se ainda não foi determinado pelo GitHub
	if pr.Mergeable == nil {
		log.Printf("⚠️ Status de merge do PR #%d ainda não foi determinado", prNumber)
		return false
	}

	// Retorna verdadeiro se não for mesclável (ou seja, tem conflito)
	return !*pr.Mergeable
}
