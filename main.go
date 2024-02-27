package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jakenvac/DotHerd/cmd"
	"github.com/jakenvac/DotHerd/config"
	"github.com/jakenvac/DotHerd/dotstore"
	"github.com/urfave/cli/v2"
)

func defaultAction(c *cli.Context) error {
	fmt.Printf("Default dot dir: %v\n", config.DEFAULT_DOT_DIR)
	return nil
}

func preRun() error {
	fileInfo, error := os.Stat(config.DEFAULT_DOT_DIR)
	noDotDir := os.IsNotExist(error)
	if !noDotDir && error != nil {
		log.Printf("Error getting path stat: %v\n", error)
		return error
	}

	if noDotDir {
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

	store, storeErr := dotstore.New()
	if storeErr != nil {
		cli.Exit(fmt.Sprintf("Error creating dotstore: %v", storeErr), 1)
	}
	defer store.Close()

	app := &cli.App{
		Version: "1.0.0-alpha",
		Commands: []*cli.Command{
			cmd.Herd(store),
			cmd.Release(store),
		},
		Name:  "All My Dots",
		Usage: "A simple dotfile manager",
		Action: func(c *cli.Context) error {
			fmt.Printf("Default dot dir: %v\n", config.DEFAULT_DOT_DIR)
			fmt.Printf("%s\n", store)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
