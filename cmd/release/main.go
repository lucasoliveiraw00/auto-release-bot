package main

import (
	"log"
	"os"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/release"
)

// Função principal
func main() {
	log.Println("🤖 Script release iniciado...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	release.ProcessReleaseEvents(*cfg)
}
