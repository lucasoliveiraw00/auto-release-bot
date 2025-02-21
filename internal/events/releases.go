package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
	date "github.com/lucasoliveiraw00/auto-release-bot/pkg/utils"
)

// Estrutura do JSON
type ReleaseEvent struct {
	Version          string `json:"version"`
	Delivery         string `json:"delivery"`
	PreparationStart string `json:"preparation_start"`
	PreparationEnd   string `json:"preparation_end"`
}

// Carrega os eventos de releases a partir de um arquivo JSON e retorna a lista de eventos junto com a data atual
func ReadReleaseEvents(cfg config.Config) ([]ReleaseEvent, time.Time, error) {
	var today time.Time
	var err error

	if cfg.MockDate != "" {
		today, err = date.ParseISODate(cfg.MockDate)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("erro ao interpretar MOCK_DATE: %v", err)
		}
		log.Printf("🔄 Usando data MOCK: %s\n", today.Format("2006-01-02"))
	} else {
		today = time.Now()
		log.Printf("📅 Data real utilizada: %s\n", today.Format("2006-01-02"))
	}

	// Leitura do arquivo JSON
	data, err := os.ReadFile(cfg.ReleaseEventsPath)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("erro ao ler o arquivo JSON: %v", err)
	}

	// Parse do JSON
	var events []ReleaseEvent
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, time.Time{}, fmt.Errorf("erro ao interpretar o JSON: %v", err)
	}

	return events, today, nil
}
