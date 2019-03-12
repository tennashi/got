package lib

import (
	"github.com/spf13/viper"
)

type GotFile struct {
	DotFile []DotFile
}

type DotFile struct {
	Dest string
	Src  string
}

func InitGotFile(path string) (*GotFile, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	gotFile := &GotFile{}
	v.Unmarshal(gotFile)
	return gotFile, nil
}
