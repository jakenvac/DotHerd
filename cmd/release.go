package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	ds "github.com/jakenvac/DotHerd/dotstore"
	"github.com/urfave/cli/v2"
)

func Release(store *ds.DotStore) *cli.Command {
	var releaseAction = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowSubcommandHelpAndExit(c, 1)
		}

		dotLink := c.Args().First()
		dotLinkAbs, absErr := filepath.Abs(dotLink)
		if absErr != nil {
			return cli.Exit(fmt.Sprintf("Error getting absolute path: %s", absErr), 1)
		}

		var aliasInPen, aliasErr = store.DotToPenAlias(dotLinkAbs)
		if aliasErr != nil {
			return cli.Exit(fmt.Sprintf("Unable to find alias in pen for dotfile %s.", dotLinkAbs), 1)
		}

		if err := os.Remove(dotLinkAbs); err != nil {
			return cli.Exit(fmt.Sprintf("Error removing symlink: %s", err), 1)
		}

		if err := os.Rename(aliasInPen, dotLinkAbs); err != nil {
			fmt.Println(err)
			return cli.Exit(fmt.Sprintf("Error moving %s to to original location: %s", aliasInPen, dotLinkAbs), 1)
		}

		if err := store.Release(dotLinkAbs); err != nil {
			return cli.Exit(fmt.Sprintf("Error removing dotfile from store: %s", err), 1)
		}

		return nil
	}

	return &cli.Command{
		Name:     "release",
		Category: "Dotfiles",
		Usage:    "Release a dotfile from the pen (Restores the original file)",
		Action:   releaseAction,
	}
}
