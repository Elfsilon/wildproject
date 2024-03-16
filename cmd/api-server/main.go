package main

import (
	"flag"
	"log"
	"temp/internal/app"
	"temp/pkg/env"

	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()

	configPath := env.String("CONFIG_PATH")
	if err := godotenv.Load(configPath); err != nil {
		log.Fatal(err)
	}

	a := app.New()
	a.Run(&app.AppFlags{
		ConfigPath: configPath,
	})
}
