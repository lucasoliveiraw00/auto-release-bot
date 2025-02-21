package prchecker

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	githubClient "github.com/google/go-github/v69/github"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/events"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/prchecker/integrations/github"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/prchecker/integrations/googlechat"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/prchecker/integrations/sonar"
	date "github.com/lucasoliveiraw00/auto-release-bot/pkg/utils"
)

// Função auxiliar para evitar valores negativos
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Função para processar eventos de release
func ProcessReleaseEvents(ctx context.Context, cfg config.Config, client *githubClient.Client) {
	releaseEvents, today, err := events.ReadReleaseEvents(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao ler eventos: %v", err)
		os.Exit(1)
	}

	for _, event := range releaseEvents {
		dates := map[string]string{
			"preparation_start": event.PreparationStart,
			"preparation_end":   event.PreparationEnd,
			"delivery":          event.Delivery,
		}

		for eventType, dateStr := range dates {
			if dateStr == "" {
				continue
			}

			eventDate, err := date.ParseISODate(dateStr)
			if err != nil {
				log.Printf("⚠️ Erro ao processar data do evento %s (%s): %v\n", event.Version, eventType, err)
				continue
			}

			if eventType == "delivery" {
				diffDays := int(eventDate.YearDay() - today.YearDay() + (eventDate.Year()-today.Year())*365)
				if diffDays == 2 || diffDays == 1 || diffDays == 0 {

					prs, _, err := client.PullRequests.List(ctx, cfg.GithubOwner, cfg.GithubRepo, nil)
					if err != nil {
						log.Fatalf("❌ Erro ao buscar PRs: %v", err)
						return
					}

					var prsInfo []map[string]string
					for _, pr := range prs {
						if pr.GetDraft() {
							continue
						}

						hasLabel := false
						for _, label := range pr.Labels {
							if label.GetName() == "RC02-02.25" {
								hasLabel = true
								break
							}
						}

						if !hasLabel {
							continue
						}

						hasSonarIssues := sonar.FetchSonarData(cfg, pr.GetNumber())
						_, failedChecks := github.CheckPRStatus(ctx, cfg, client, pr.GetHead().GetSHA())
						hasPendingChanges, hasUnresolvedComments, approvalCount := github.GetApprovalStatus(ctx, cfg, client, pr.GetNumber())
						hasConflicts := github.HasMergeConflicts(ctx, cfg, client, pr.GetNumber())

						pendingApprovals := max(0, 2-approvalCount)
						if failedChecks > 0 || hasPendingChanges || hasUnresolvedComments || hasSonarIssues || pendingApprovals > 0 || hasConflicts {
							details := ""

							if pendingApprovals != 0 {
								details += fmt.Sprintf("✅ Aprovações: %d/2 - Falta %d aprovação(s)", approvalCount, pendingApprovals)
							}
							if failedChecks > 0 {
								details += fmt.Sprintf("<br>🔍 Checks: %d falha(s)", failedChecks)
							}
							if hasPendingChanges {
								details += "<br>🚨 Possui mudanças pendentes"
							}
							if hasUnresolvedComments {
								details += "<br>💬 Possui comentários não resolvidos"
							}
							if hasConflicts {
								details += "<br>⚠️ Este PR tem conflitos e precisa ser resolvido!"
							}
							if hasSonarIssues {
								details += "<br>❌ O SonarCloud detectou pendências de qualidade"
							}

							title := pr.GetTitle()
							if len(title) > 72 {
								title = title[:72] + "..."
							}

							prsInfo = append(prsInfo, map[string]string{
								"number":  fmt.Sprintf("🔹 <b>PR #%s</b>", strconv.Itoa(pr.GetNumber())),
								"title":   fmt.Sprintf("⚓️ <a href='%s'>%s</a>", pr.GetHTMLURL(), title),
								"url":     pr.GetHTMLURL(),
								"details": details,
							})
						}
					}

					if len(prsInfo) == 0 {
						log.Println("✅ Nenhum PR com pendências. Nenhuma notificação enviada.")
						return
					}

					title := fmt.Sprintf("<br><b>⏳ A data de entrega da %s está chegando!</b><br> Fique de olho no tempo restante para garantir que haja tempo suficiente para corrigir os PRs com pendências.<br><br>", event.Version)

					if diffDays == 0 {
						title = fmt.Sprintf("<br><b>📅 Chegou o grande dia da entrega da %s!</b><br><font color=\"#ffc107\">Vamos dar aquele último gás para deixar tudo certinho. Reserve um tempo para corrigir os PRs que ainda têm pendências.</font><br><br>", event.Version)
					}

					prsInfo[len(prsInfo)-1]["details"] += "<br>"

					googlechat.SendToGoogleChat(cfg, prsInfo, title, event.Version)
					return
				}
			}
		}
	}

	log.Println("✅ Nenhum PR com pendências. Nenhuma notificação enviada.")
}
