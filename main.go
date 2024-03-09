package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jakenvac/DotHerd/cmd"
	"github.com/jakenvac/DotHerd/config"
	"github.com/jakenvac/DotHerd/repo"
	"github.com/urfave/cli/v2"
)

func preRun() error {
	fileInfo, error := os.Stat(config.DEFAULT_DOT_DIR)
	dotDirNotExist := os.IsNotExist(error)
	if !dotDirNotExist && error != nil {
		log.Printf("Error getting path stat: %v\n", error)
		return error
	}

	if dotDirNotExist {
		err := os.MkdirAll(config.DEFAULT_DOT_DIR, 0755)
		if err != nil {
			log.Printf("Error creating dot dir: %v\n", err)
			return err
		}
		log.Printf("Created dot dir: %v\n", config.DEFAULT_DOT_DIR)
	}

	if fileInfo != nil && !fileInfo.IsDir() {
		log.Printf("Error: %v is not a directory\n", config.DEFAULT_DOT_DIR)
		return fmt.Errorf("Error: %v is not a directory", config.DEFAULT_DOT_DIR)
	}

	return nil
}

func main() {
	preErr := preRun()
	if preErr != nil {
		log.Fatal(fmt.Sprintf("Error during Dot Dir setup: %v", preErr))
	}

	repo, repoErr := repo.New()
	if repoErr != nil {
		cli.Exit(fmt.Sprintf("Error creating dot repo: %v", repoErr), 1)
	}
	defer repo.Close()

	app := &cli.App{
		Version: "1.0.0-alpha",
		Commands: []*cli.Command{
			cmd.Herd(repo),
			cmd.Unherd(repo),
			cmd.Json(repo),
		},
		Name:  "DotHerd",
		Usage: "Herd your dotfiles into one place",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
