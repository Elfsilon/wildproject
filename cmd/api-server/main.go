package main

import (
	"flag"
	"log"
	"temp/internal/app"

	"github.com/joho/godotenv"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "cfg", "", "location of the config.yaml file")
}

func main() {
	flag.Parse()

	if err := godotenv.Load(configPath); err != nil {
		log.Fatal(err)
	}

	a := app.New()
	a.Run(&app.AppFlags{
		ConfigPath: configPath,
	})
}
