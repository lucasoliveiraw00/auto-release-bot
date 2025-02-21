package main

import (
	"context"
	"log"
	"os"

	"github.com/lucasoliveiraw00/auto-release-bot/config"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/prchecker"
	"github.com/lucasoliveiraw00/auto-release-bot/internal/prchecker/integrations/github"
)

// Função principal
func main() {
	log.Println("🤖 Script prchecker iniciado...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ctx := context.Background()
	client := github.NewClient(ctx, *cfg)

	prchecker.ProcessReleaseEvents(ctx, *cfg, client)
}
