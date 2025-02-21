package release

import (
	"log"
	"os"
	"time"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/events"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/release/integrations/googlechat"
)

// ProcessReleaseEvents verifica os eventos e envia notificações
func ProcessReleaseEvents(cfg config.Config) {
	releaseEvents, today, err := events.ReadReleaseEvents(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao ler eventos: %v", err)
		os.Exit(1)
	}
	for _, event := range releaseEvents {
		eventDates := map[string]string{
			"preparation_start": event.PreparationStart,
			"preparation_end":   event.PreparationEnd,
			"delivery":          event.Delivery,
		}

		for eventType, dateStr := range eventDates {
			if dateStr == "" {
				continue
			}

			eventDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				log.Printf("⚠️ Erro ao processar data do evento %s (%s): %v\n", event.Version, eventType, err)
				continue
			}

			diffDays := int(eventDate.YearDay() - today.YearDay() + (eventDate.Year()-today.Year())*365)

			var alertType string
			switch {
			case eventType == "preparation_start" && diffDays == 0:
				alertType = "preparation_start"
			case eventType == "preparation_end" && diffDays == 0:
				alertType = "preparation_end"
			case eventType == "delivery" && diffDays == 2:
				alertType = "delivery_2days"
			case eventType == "delivery" && diffDays == 1:
				alertType = "delivery_1day"
			case eventType == "delivery" && diffDays == 0:
				alertType = "delivery_today"
			}

			if alertType != "" {
				googlechat.SendToGoogleChat(cfg, alertType, event.Version)
				log.Printf("✅ Notificação enviada para '%s' da versão %s\n", alertType, event.Version)
				return
			}
		}
	}

	log.Println("✅ Nenhuma notificação necessária.")
}
