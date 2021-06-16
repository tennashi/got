package cli

import (
	"fmt"
	"os"

	"github.com/tennashi/got"
	"github.com/urfave/cli/v2"
)

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list installed packages",
		Flags: []cli.Flag{
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

			if !c.IsSet("datadir") {
				if err := c.Set("datadir", cf.DataDir); err != nil {
					return err
				}
			}

			if c.NArg() != 0 {
				return &got.InvalidParamError{
					Message: "no params required",
				}
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			cfg := &got.ListCommandConfig{
				DataDir: c.Path("datadir"),
				IsDebug: c.Bool("debug"),
			}

			ioStream := &got.IOStream{
				In:  c.App.Reader,
				Out: c.App.Writer,
				Err: c.App.ErrWriter,
			}

			cmd, err := got.NewListCommand(ioStream, cfg)
			if err != nil {
				return err
			}

			return cmd.Run()
		},
	}
}
