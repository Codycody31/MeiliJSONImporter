package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"go.codycody31.dev/MeiliJSONImporter/pkg/importer"
)

func main() {
	app := &cli.App{
		Name:   "MeiliJSONImporter",
		Usage:  "Import JSON data into MeiliSearch",
		Flags:  importer.SetupFlags(),
		Action: importer.HandleImport,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
