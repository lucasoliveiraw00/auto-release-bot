package sonar

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
)

type SonarResponse struct {
	Component struct {
		Measures []struct {
			Metric  string `json:"metric"`
			Periods []struct {
				Value string `json:"value"`
			} `json:"periods"`
		} `json:"measures"`
	} `json:"component"`
}

// Verifica se o PR tem pendências no SonarCloud
func FetchSonarData(cfg config.Config, prNumber int) bool {
	url := fmt.Sprintf("https://sonarcloud.io/api/measures/component?pullRequest=%d&metricKeys=new_violations,new_coverage,new_duplicated_lines_density,new_security_hotspots,new_accepted_issues&component=%s", prNumber, cfg.SonarOwner)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("⚠️ Erro ao criar requisição para Sonar: %v", err)
		return false
	}

	req.SetBasicAuth(cfg.SonarToken, "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar dados do Sonar: %v", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("⚠️ Erro ao ler resposta do Sonar: %v", err)
		return false
	}

	var sonarData SonarResponse
	if err := json.Unmarshal(body, &sonarData); err != nil {
		log.Printf("⚠️ Erro ao decodificar JSON do Sonar: %v", err)
		return false
	}

	thresholds := map[string]float64{
		"new_security_hotspots":        cfg.SonarThresholds.NewSecurityHotspots,
		"new_violations":               cfg.SonarThresholds.NewViolations,
		"new_accepted_issues":          cfg.SonarThresholds.NewAcceptedIssues,
		"new_coverage":                 cfg.SonarThresholds.NewCoverage,
		"new_duplicated_lines_density": cfg.SonarThresholds.NewDuplicatedLinesDensity,
	}

	issues := make(map[string]float64)

	for _, measure := range sonarData.Component.Measures {
		if len(measure.Periods) > 0 {
			value, _ := strconv.ParseFloat(measure.Periods[0].Value, 64)
			if threshold, exists := thresholds[measure.Metric]; exists {
				if (measure.Metric == "new_coverage" && value < threshold) ||
					(measure.Metric != "new_coverage" && value > threshold) {
					issues[measure.Metric] = value
					log.Printf("⚠️ Problema detectado no Sonar: %s = %f (esperado ≤ %f)", measure.Metric, value, threshold)
				}
			}
		}
	}

	if len(issues) == 0 {
		log.Println("✅ Nenhum problema encontrado no SonarCloud!")
	}

	return len(issues) > 0
}
