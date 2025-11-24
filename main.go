package main

import (
	"log"

	"github.com/axellelanca/urlshortener/cmd"
	_ "github.com/axellelanca/urlshortener/cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "github.com/axellelanca/urlshortener/cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
	"github.com/axellelanca/urlshortener/internal/config"
)

func main() {
	// TODO
	//Boris : Chargement de la configuration
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cmd.Execute()

}
