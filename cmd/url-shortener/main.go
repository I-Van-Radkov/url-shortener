package main

import (
	"log"

	"github.com/I-Van-Radkov/url-shortener/internal/app"
	"github.com/I-Van-Radkov/url-shortener/internal/config"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("app init error: %v", err)
	}

	if err = a.Run(); err != nil {
		log.Fatalf("app run error: %v", err)
	}
}
