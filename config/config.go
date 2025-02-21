package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
}

// LoadConfig carrega as variáveis de ambiente
func LoadConfig() (*Config, error) {
	envPath := flag.String("env", "./.env", "Caminho para o arquivo .env (padrão: ./.env)")
	releaseEventsPath := flag.String("release-events", "./config/release_events.json", "Caminho para o arquivo release_events.json (padrão: ./config/release_events.json)")
	flag.Parse()

	err := godotenv.Load(*envPath)
	if err != nil {
		log.Println("⚠️ Nenhum arquivo .env foi encontrado, usando variáveis de ambiente padrão")
	}

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

	// Verifica se todas as variáveis obrigatórias estão preenchidas
	for _, envVar := range requiredEnvVars {
		if envVar.value == "" {
			return nil, fmt.Errorf("❌ ERRO: Variável obrigatória ausente no ambiente: %s", envVar.name)
		}
	}

	return config, nil
}

// getEnv pega variável de ambiente com fallback para valor padrão
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
