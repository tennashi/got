package cli

import (
	"fmt"
	"os"

	"github.com/tennashi/got"
	"github.com/urfave/cli/v2"
)

func NewDisableCommand() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "disable the specified executable",
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

			if c.NArg() != 2 {
				return &got.InvalidParamError{
					Message: "package name and executable name is required",
				}
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			cfg := &got.DisableCommandConfig{
				DataDir: c.Path("datadir"),
				BinDir:  c.Path("bindir"),
				IsDebug: c.Bool("debug"),
			}

			ioStream := &got.IOStream{
				In:  c.App.Reader,
				Out: c.App.Writer,
				Err: c.App.ErrWriter,
			}

			cmd, err := got.NewDisableCommand(ioStream, cfg)
			if err != nil {
				return err
			}

			pkgName := c.Args().Get(0)
			execName := c.Args().Get(1)
			return cmd.Run(pkgName, execName)
		},
	}
}
