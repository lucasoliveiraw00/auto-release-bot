package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// SonarThresholds armazena os critérios do Sonar
type SonarThresholds struct {
	NewSecurityHotspots       float64
	NewViolations             float64
	NewAcceptedIssues         float64
	NewCoverage               float64
	NewDuplicatedLinesDensity float64
}

// Config estrutura todas as variáveis de ambiente da aplicação
type Config struct {
	GoogleChatWebhook  string
	MockDate           string
	GithubOwner        string
	GithubRepo         string
	GithubToken        string
	SonarToken         string
	SonarOwner         string
	VersionCalendarURL string
	QualityCriteriaURL string
	ReleaseEventsPath  string
	SonarThresholds    SonarThresholds
}

// LoadConfig carrega todas as variáveis do .env e do sistema
func LoadConfig() (*Config, error) {
	envPath := flag.String("env", "./.env", "Caminho para o arquivo .env (padrão: ./.env)")
	releaseEventsPath := flag.String("release-events", "./config/release_events.json", "Caminho para o arquivo release_events.json (padrão: ./config/release_events.json)")
	flag.Parse()

	err := godotenv.Load(*envPath)
	if err != nil {
		log.Println("⚠️ Nenhum arquivo .env foi encontrado, usando variáveis de ambiente padrão")
	}

	// Lendo e convertendo os thresholds do Sonar
	sonarThresholds := SonarThresholds{
		NewSecurityHotspots:       getEnvAsFloat("SONAR_NEW_SECURITY_HOTSPOTS", 0),
		NewViolations:             getEnvAsFloat("SONAR_NEW_VIOLATIONS", 0),
		NewAcceptedIssues:         getEnvAsFloat("SONAR_NEW_ACCEPTED_ISSUES", 0),
		NewCoverage:               getEnvAsFloat("SONAR_NEW_COVERAGE", 75),
		NewDuplicatedLinesDensity: getEnvAsFloat("SONAR_NEW_DUPLICATED_LINES_DENSITY", 0),
	}

	// Monta a struct Config
	config := &Config{
		GoogleChatWebhook:  getEnv("GOOGLE_CHAT_WEBHOOK", ""),
		MockDate:           getEnv("MOCK_DATE", ""),
		GithubOwner:        getEnv("GITHUB_OWNER", ""),
		GithubRepo:         getEnv("GITHUB_REPO", ""),
		GithubToken:        getEnv("GITHUB_TOKEN", ""),
		SonarToken:         getEnv("SONAR_TOKEN", ""),
		SonarOwner:         getEnv("SONAR_OWNER", ""),
		VersionCalendarURL: getEnv("VERSION_CALENDAR_URL", ""),
		QualityCriteriaURL: getEnv("QUALITY_CRITERIA_URL", ""),
		ReleaseEventsPath:  *releaseEventsPath,
		SonarThresholds:    sonarThresholds,
	}

	// Lista de variáveis obrigatórias
	requiredEnvVars := []struct {
		name  string
		value string
	}{
		{"GITHUB_OWNER", config.GithubOwner},
		{"GITHUB_REPO", config.GithubRepo},
		{"GITHUB_TOKEN", config.GithubToken},
		{"SONAR_TOKEN", config.SonarToken},
		{"SONAR_OWNER", config.SonarOwner},
		{"GOOGLE_CHAT_WEBHOOK", config.GoogleChatWebhook},
		{"VERSION_CALENDAR_URL", config.VersionCalendarURL},
		{"QUALITY_CRITERIA_URL", config.QualityCriteriaURL},
	}

	// Valida as variáveis obrigatórias
	for _, envVar := range requiredEnvVars {
		if envVar.value == "" {
			return nil, fmt.Errorf("❌ ERRO: Variável obrigatória ausente no ambiente: %s", envVar.name)
		}
	}

	return config, nil
}

// getEnv lê uma variável de ambiente com fallback para um valor padrão
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsFloat lê uma variável de ambiente e a converte para float64
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	valFloat, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		log.Printf("⚠️ Erro ao converter %s='%s' para float64. Usando valor padrão %.2f", key, valStr, defaultValue)
		return defaultValue
	}
	return valFloat
}
