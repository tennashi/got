package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func Run() int {
	app := &cli.App{
		Name:  "got",
		Usage: "package manager for commands written in Go",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Usage:   "config file",
				Aliases: []string{"c"},
				EnvVars: []string{"GOT_CONFIG_FILE"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "debug mode",
				Value:   false,
				EnvVars: []string{"GOT_DEBUG"},
			},
		},
		Commands: []*cli.Command{
			NewInstallCommand(),
			NewUpgradeCommand(),
			NewListCommand(),
			NewEnableCommand(),
			NewDisableCommand(),
			NewShowCommand(),
			NewPinCommand(),
			NewUnpinCommand(),
			NewRemoveCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.ErrWriter, "an error occurred during command running: %v\n", err)
		return 1
	}

	return 0
}
