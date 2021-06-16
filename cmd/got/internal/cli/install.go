package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tennashi/got"
	"github.com/urfave/cli/v2"
)

func NewInstallCommand() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "install the specified package",
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

			if !c.IsSet("all") {
				if err := c.Set("all", strconv.FormatBool(cf.InstallAll)); err != nil {
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
			cfg := &got.InstallCommandConfig{
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

			cmd, err := got.NewInstallCommand(ioStream, cfg)
			if err != nil {
				return err
			}

			pkgName := c.Args().Get(0)
			return cmd.Run(pkgName)
		},
	}
}
