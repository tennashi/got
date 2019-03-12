package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tennashi/got/lib"
)

func NewSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync [remote repository [local directory]]",
		Short: "clone your dotfiles repository",
		Run:   doSync,
	}
}

func doSync(c *cobra.Command, args []string) {
	var remoteURL string
	var localDir string

	switch len(args) {
	case 0:
		if globalConfig.ConfigFileUsed() == "" {
			c.Println("Config file unloaded.")
			c.Println("Please\n\t'got -c /path/to/config.toml sync'\nor\n\t'got sync https://your.remote.repository/dotfiles.git /path/to/dotfiles'.\n")
			c.Usage()
			return
		}
		remoteURL = globalConfig.GetString("dotfiles.remote")
		localDir = globalConfig.GetString("dotfiles.local")
	case 1:
		remoteURL = args[0]
		localDir = globalConfig.GetString("dotfiles.local")
	case 2:
		remoteURL = args[0]
		localDir = args[1]
	default:
		c.Usage()
		return
	}

	localDir, err := lib.ExpandPath(localDir)
	if err != nil {
		c.Println(err)
		return
	}

	out, _ := exec.Command("git", "clone", remoteURL, localDir).CombinedOutput()
	c.Printf("git: %s", string(out))

	gotFilePath := filepath.Join(localDir, "Gotfile.toml")
	gotFile, err := lib.InitGotFile(gotFilePath)
	if err != nil {
		return
	}
	for _, dotFile := range gotFile.DotFile {
		srcPath := filepath.Join(localDir, dotFile.Src)
		destPath, err := lib.ExpandPath(dotFile.Dest)
		if err != nil {
			c.Println(err)
			return
		}

		if _, err := os.Stat(destPath); err == nil {
			c.Println("exist dest file:", destPath)
			continue
		}

		if info, err := os.Stat(srcPath); err != nil {
			c.Println("not exist src file:", srcPath)
			continue
		} else if info.IsDir() {
			srcPath += "/"
		}
		c.Println("ln -s", srcPath, destPath)
		out, _ = exec.Command("ln", "-s", srcPath, destPath).CombinedOutput()
		c.Printf("ln: %s", string(out))
	}
}

func init() {
	rootCmd.AddCommand(NewSyncCmd())
}
