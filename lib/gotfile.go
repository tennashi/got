package lib

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type GotFile struct {
	DotFile []DotFile
}

type DotFile struct {
	Dest string
	Src  string
}

func InitGotFile(dirPath string) (*GotFile, error) {
	path := filepath.Join(dirPath, "Gotfile.toml")
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	gotFile := &GotFile{}
	viper.Unmarshal(gotFile)
	return gotFile, nil
}
