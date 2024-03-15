package main

import (
	"flag"
	"temp/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "cfg", "", "location of the config.yaml file")
}

func main() {
	flag.Parse()

	a := app.New()
	a.Run(&app.AppFlags{
		ConfigPath: configPath,
	})
}
