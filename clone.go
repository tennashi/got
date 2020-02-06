package got

import (
	"github.com/tennashi/got/io"
	"github.com/tennashi/got/repository"
)

func clone(ioStream *io.Stream, repoName string) error {
	repo, err := repository.New(ioStream, repoName)
	if err != nil {
		return err
	}
	return repo.Clone()
}
