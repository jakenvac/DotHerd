package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jakenvac/DotHerd/repo"
	"github.com/urfave/cli/v2"
)

func Unherd(repo *repo.DotRepo) *cli.Command {
	var unherdAction = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowSubcommandHelpAndExit(c, 1)
		}

		dotLink := c.Args().First()
		dotLinkAbs, absErr := filepath.Abs(dotLink)
		if absErr != nil {
			return cli.Exit(fmt.Sprintf("Error getting absolute path: %s", absErr), 1)
		}

		var nameInPen, nameErr = repo.NameFromLink(dotLinkAbs)
		if nameErr != nil {
			return cli.Exit(fmt.Sprintf("Unable to find alias in pen for dotfile %s.", dotLinkAbs), 1)
		}

		if err := os.Remove(dotLinkAbs); err != nil {
			return cli.Exit(fmt.Sprintf("Error removing symlink: %s", err), 1)
		}

		if err := os.Rename(nameInPen, dotLinkAbs); err != nil {
			fmt.Println(err)
			return cli.Exit(fmt.Sprintf("Error moving %s to to original location: %s", nameInPen, dotLinkAbs), 1)
		}

		if err := repo.Unherd(dotLinkAbs); err != nil {
			return cli.Exit(fmt.Sprintf("Error removing dotfile from repo: %s", err), 1)
		}

		return nil
	}

	return &cli.Command{
		Name:     "unherd",
		Category: "Dotfiles",
		Usage:    "Unherd a dotfile from the pen (Rerepos the original file)",
		Action:   unherdAction,
	}
}
