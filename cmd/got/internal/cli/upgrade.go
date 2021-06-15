package cli

import (
	"errors"
	"strconv"

	"github.com/tennashi/got"
	"github.com/urfave/cli/v2"
)

func NewUpgradeCommand() *cli.Command {
	return &cli.Command{
		Name:  "upgrade",
		Usage: "upgrade installed packages",
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
			&cli.BoolFlag{
				Name:    "all",
				Usage:   "install all command",
				Aliases: []string{"a"},
				EnvVars: []string{"GOT_INSTALL_ALL"},
			},
		},
		Before: func(c *cli.Context) error {
			cf, err := NewConfigFile(c.Path("config"))
			if err != nil {
				return err
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

			if !c.IsSet("all") {
				if err := c.Set("all", strconv.FormatBool(cf.InstallAll)); err != nil {
					return err
				}
			}

			if c.NArg() != 0 {
				return errors.New("some error")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			cfg := &got.UpgradeCommandConfig{
				DataDir:           c.Path("datadir"),
				BinDir:            c.Path("bindir"),
				InstallAllCommand: c.Bool("all"),
				IsDebug:           c.Bool("debug"),
			}

			ioStream := &got.IOStream{
				In:  c.App.Reader,
				Out: c.App.Writer,
				Err: c.App.ErrWriter,
			}

			cmd, err := got.NewUpgradeCommand(ioStream, cfg)
			if err != nil {
				return err
			}

			return cmd.Run()
		},
	}
}
