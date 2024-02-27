package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/jakenvac/DotHerd/config"
	ds "github.com/jakenvac/DotHerd/dotstore"
	"github.com/urfave/cli/v2"
)

func Herd(store *ds.DotStore) *cli.Command {
	var herdAction = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowSubcommandHelpAndExit(c, 1)
		}

		// @TODO replace home directory with {{home}} in path
		dot := c.Args().First()
		absDot, absErr := filepath.Abs(dot)
		if absErr != nil {
			return cli.Exit(fmt.Sprintf("Error getting absolute path: %s", absErr), 1)
		}

		var destinationName string
		if alias := c.String("alias"); alias != "" {
			destinationName = fmt.Sprintf("%s%s", alias, path.Ext(absDot))
		} else {
			destinationName = absDot
		}

		destinationInPen := path.Join(config.DEFAULT_DOT_DIR, destinationName)
		if _, err := os.Stat(destinationInPen); err == nil {
			return cli.Exit(fmt.Sprintf("Dotfile %s already exists in Dot Dir", destinationName), 1)
		}

		if err := os.Rename(absDot, destinationInPen); err != nil {
			return cli.Exit(fmt.Sprintf("Error moving %s to Dot Dir", err), 1)
		}

		if err := os.Symlink(destinationInPen, absDot); err != nil {
			return cli.Exit(fmt.Sprintf("Error creating symlink: %s", err), 1)
		}

		if err := store.Herd(absDot, destinationInPen); err != nil {
			return cli.Exit(fmt.Sprintf("Error adding dotfile to store: %s", err), 1)
		}

		fmt.Printf("Dotfile %s tracked as %s\n", absDot, destinationName)

		return nil
	}

	var aliasFlag = &cli.StringFlag{
		Name:    "alias",
		Usage:   "Alias for the dotfile, if not provided the filename will be used",
		Aliases: []string{"a"},
	}

	return &cli.Command{
		Name:     "herd",
		Category: "Dotfiles",
		Usage:    "Herd a new dotfile or directory into the pen",
		Action:   herdAction,
		Flags: []cli.Flag{
			aliasFlag,
		},
	}
}
