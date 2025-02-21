package googlechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
)

// Envia mensagem para o Google Chat no formato correto
func SendToGoogleChat(cfg config.Config, prsInfo []map[string]string, title, release string) {
	var widgets []map[string]interface{}
	widgets = append(widgets, map[string]interface{}{
		"textParagraph": map[string]interface{}{
			"text": title,
		},
	})

	for i, pr := range prsInfo {
		widgets = append(widgets, map[string]interface{}{
			"textParagraph": map[string]interface{}{
				"text": fmt.Sprintf("%s<users/all> <br>%s<br>%s<br>",
					pr["number"], pr["title"], pr["details"]),
			},
		})
		if i != len(prsInfo)-1 {
			widgets = append(widgets, map[string]interface{}{
				"divider": map[string]interface{}{},
			})
		}
	}

	widgets = append(widgets, map[string]interface{}{
		"buttonList": map[string]interface{}{
			"buttons": []map[string]interface{}{
				{
					"text": "Ver critérios de qualidade",
					"icon": map[string]interface{}{
						"materialIcon": map[string]interface{}{
							"name": "search",
						},
					},
					"type": "OUTLINED",
					"onClick": map[string]interface{}{
						"openLink": map[string]interface{}{
							"url": cfg.QualityCriteriaURL,
						},
					},
				},
				{
					"text": "Ver calendário de versões",
					"icon": map[string]interface{}{
						"materialIcon": map[string]interface{}{
							"name": "calendar_month",
						},
					},
					"type": "OUTLINED",
					"onClick": map[string]interface{}{
						"openLink": map[string]interface{}{
							"url": cfg.VersionCalendarURL,
						},
					},
				},
			},
		},
	})

	payload := map[string]interface{}{
		"cardsV2": []map[string]interface{}{
			{
				"cardId": "notificacao_card",
				"card": map[string]interface{}{
					"header": map[string]interface{}{
						"title":    fmt.Sprintf("📦 Pull Requests: %s", release),
						"subtitle": "PRs que possuem pendências",
					},
					"sections": []map[string]interface{}{
						{
							"collapsible":               false,
							"uncollapsibleWidgetsCount": 1,
							"widgets":                   widgets,
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("⚠️ Erro ao serializar JSON: %v", err)
		return
	}

	resp, err := http.Post(cfg.GoogleChatWebhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("⚠️ Erro ao enviar mensagem para Google Chat: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("✅ Mensagem enviada com sucesso para o Google Chat!")
	} else {
		log.Printf("⚠️ Falha no envio ao Google Chat, status: %s", resp.Status)
	}
}
