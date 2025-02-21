package googlechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
)

// Gera o payload estruturado no formato de card
func generatePayload(cfg config.Config, eventType, version string) map[string]interface{} {
	description := getMessage(cfg, eventType, version)
	if description == "" {
		return nil
	}

	return map[string]interface{}{
		"cardsV2": []map[string]interface{}{
			{
				"cardId": "notificacao_card",
				"card": map[string]interface{}{
					"header": map[string]interface{}{
						"subtitle": "Calendário de versões",
						"title":    fmt.Sprintf("📦 <b>%s</b>", version),
					},
					"sections": []map[string]interface{}{
						{
							"header":                    "🔔 Atualização",
							"collapsible":               false,
							"uncollapsibleWidgetsCount": 1,
							"widgets": []map[string]interface{}{
								{
									"textParagraph": map[string]interface{}{
										"text": fmt.Sprintf("%s\n\n", description),
									},
								},
								{
									"buttonList": map[string]interface{}{
										"buttons": []map[string]interface{}{
											{
												"text": "Ver calendário de versões",
												"onClick": map[string]interface{}{
													"openLink": map[string]interface{}{
														"url": cfg.VersionCalendarURL,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// getMessage retorna o título e a mensagem correspondente ao evento
func getMessage(cfg config.Config, eventType, version string) string {
	var description string

	switch eventType {
	case "preparation_start":
		description = fmt.Sprintf(
			"<users/all> <br><b>🎯 Hoje iniciamos a preparação para %s!</b><br><br>"+
				"Entramos na fase de testes das demandas alocadas na RC! 📦<br>"+
				"É o momento de validar cada detalhe e garantir que tudo esteja pronto para o lançamento. 🚀<br><br>"+
				"Vamos juntos garantir a qualidade e o sucesso dessa entrega! 🚀", version)
	case "preparation_end":
		description = fmt.Sprintf(
			"<users/all> <br><b>🎯 Último dia de preparação para %s!</b><br><br>"+
				"Chegamos à reta final da fase de testes! 🌟<br>"+
				"Agora é o momento de revisar e garantir que todas as demandas alocadas na RC estejam validadas com sucesso. 🔍<br><br>"+
				"Vamos juntos garantir a qualidade e o sucesso dessa entrega! 🚀", version)
	case "delivery_2days":
		description = fmt.Sprintf(
			"<users/all> <br><b>📢 Faltam apenas 2 dias para a entrega da %s!</b><br><br>"+
				"O prazo final está se aproximando! ⏳<br>"+
				"Agora é o momento de revisar todos os detalhes, validar se as demandas atendem aos requisitos de qualidade. ⚙️<br><br>"+
				"Se houver algo pendente, esta é a melhor oportunidade para ajustes finais! 🚀", version)
	case "delivery_1day":
		description = fmt.Sprintf(
			"<users/all> <br><b><font color=\"#ffc107\">⏳ Amanhã é o prazo final para entrega da %s!</font></b><br><br>"+
				"Estamos a apenas um dia da entrega! 📆<br>"+
				"Garanta que todas as demandas estejam devidamente finalizadas e revisadas. Se houver pendências, este é o último momento para resolvê-las! 🛠<br><br>"+
				"Vamos juntos garantir uma entrega impecável! 🚀", version)
	case "delivery_today":
		description = fmt.Sprintf(
			"<users/all> <br><b><font color=\"#ffc107\">📅 Hoje é o último dia para entrega das demandas da %s!</font></b><br><br>"+
				"Se você tem alguma demanda devidamente pronta, seguindo todos os <a href=\"%s\">critérios de qualidade</a>, ainda há tempo para alocá-la na release! 🚀<br><br>"+
				"Caso ainda falte algum ajuste, aproveite este momento para finalizar as correções necessárias. "+
				"<font color=\"#ffc107\">Lembrando que apenas demandas completas e em conformidade com os "+
				"<a href=\"%s\">critérios de qualidade</a> poderão ser incluídas na RC.</font></b><br><br>"+
				"Vamos juntos garantir uma release impecável! 🧡 ✨",
			version, cfg.QualityCriteriaURL, cfg.QualityCriteriaURL)
	default:
		return ""
	}

	return description
}

// SendToGoogleChat envia notificações para o Google Chat
func SendToGoogleChat(cfg config.Config, eventType, version string) {
	payload := generatePayload(cfg, eventType, version)
	if payload == nil {
		log.Println("⚠️ Nenhum payload gerado, não enviando notificação.")
		return
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("⚠️ Erro ao serializar JSON: %v", err)
		return
	}

	resp, err := http.Post(cfg.GoogleChatWebhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("⚠️ Erro ao enviar requisição para Google Chat: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("✅ Mensagem enviada com sucesso para o Google Chat!")
	} else {
		log.Printf("⚠️ Falha no envio ao Google Chat, status: %s", resp.Status)
	}
}
