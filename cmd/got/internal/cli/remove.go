package cli

import (
	"fmt"
	"os"

	"github.com/tennashi/got"
	"github.com/urfave/cli/v2"
)

func NewRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove the specified package",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "bindir",
				Usage:   "bin directory",
				Aliases: []string{"b"},
				EnvVars: []string{"GOT_BIN_DIR"},
			},
			&cli.PathFlag{
				Name:    "datadir",
				Usage:   "data directory",
				Aliases: []string{"d"},
				EnvVars: []string{"GOT_DATA_DIR"},
			},
		},
		Before: func(c *cli.Context) error {
			cf, err := NewConfigFile(c.Path("config"))
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}

				fmt.Fprintln(c.App.ErrWriter, "Load the default config")
			}

			if !c.IsSet("bindir") {
				if err := c.Set("bindir", cf.BinDir); err != nil {
					return err
				}
			}

			if !c.IsSet("datadir") {
				if err := c.Set("datadir", cf.DataDir); err != nil {
					return err
				}
			}

			if c.NArg() != 1 {
				return &got.InvalidParamError{
					Message: "package name is required",
				}
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			cfg := &got.RemoveCommandConfig{
				DataDir: c.Path("datadir"),
				BinDir:  c.Path("bindir"),
				IsDebug: c.Bool("debug"),
			}

			ioStream := &got.IOStream{
				In:  c.App.Reader,
				Out: c.App.Writer,
				Err: c.App.ErrWriter,
			}

			cmd, err := got.NewRemoveCommand(ioStream, cfg)
			if err != nil {
				return err
			}

			pkgName := c.Args().Get(0)
			return cmd.Run(pkgName)
		},
	}
}
