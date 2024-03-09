package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/jakenvac/DotHerd/config"
	"github.com/jakenvac/DotHerd/repo"
	"github.com/urfave/cli/v2"
)

func herdAction(repo *repo.DotRepo, c *cli.Context) error {
	if c.NArg() == 0 || c.NArg() > 2 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	dot := c.Args().First()
	absDot, absErr := filepath.Abs(dot)
	if absErr != nil {
		return cli.Exit(fmt.Sprintf("Error getting absolute path: %s", absErr), 1)
	}

	var destinationName string
	var fileName = path.Base(dot)

	if nameInPen := c.Args().Get(1); nameInPen != "" {
		destinationName = nameInPen 
	} else {
		destinationName = fileName
	}

	destinationInPen := path.Join(config.DEFAULT_DOT_DIR, destinationName)
	if _, err := os.Stat(destinationInPen); err == nil {
		return cli.Exit(fmt.Sprintf("Dotfile %s already exists in Dot Dir", destinationName), 1)
	}

	if err := os.Rename(absDot, destinationInPen); err != nil {
		return cli.Exit(fmt.Sprintf("Error moving %s to Dot Dir: %s", dot, err), 1)
	}

	if err := os.Symlink(destinationInPen, absDot); err != nil {
		return cli.Exit(fmt.Sprintf("Error creating symlink: %s", err), 1)
	}

	// @eODO replace home directory with {{home}} in path
	// @TODO replace dotdir with {{dotdir}} in path
	if err := repo.Herd(absDot, destinationInPen); err != nil {
		return cli.Exit(fmt.Sprintf("Error adding dotfile to store: %s", err), 1)
	}

	fmt.Printf("Dotfile %s tracked as %s\n", absDot, destinationName)

	return nil
}

func Herd(repo *repo.DotRepo) *cli.Command {
	return &cli.Command{
		Name:     "herd",
		Description: "Add a new dotfile to the pen",
		Usage:    "dot herd <dotfile> [name in pen]",
		Action: func(c *cli.Context) error {
			return herdAction(repo, c)
		},
	}
}
