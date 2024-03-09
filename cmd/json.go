package cmd

import (
	"fmt"

	"github.com/jakenvac/DotHerd/repo"
	"github.com/urfave/cli/v2"
)

func Json(repo *repo.DotRepo) *cli.Command {
	return &cli.Command{
		Name:  "json",
		Usage: "Output the dot repo as JSON for source control",
		Action: func(c *cli.Context) error {
			if json, err := repo.Json(); err != nil {
				return err
			} else {
				fmt.Println(json)
			}

			return nil
		},
	}
}
