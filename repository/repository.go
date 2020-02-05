package repository

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	app_io "github.com/tennashi/got/io"
)

type Repository struct {
	URI string

	ioStream *app_io.Stream
}

func New(ioStream *app_io.Stream, repoName string) (*Repository, error) {
	uri, err := parseRepoName(repoName)
	if err != nil {
		return nil, err
	}
	return &Repository{
		URI:      uri,
		ioStream: ioStream,
	}, nil
}

func (r *Repository) Cmd(args []string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Stdout = r.ioStream.Out
	cmd.Stderr = r.ioStream.Err
	cmd.Stdin = r.ioStream.In
	return cmd
}

func parseRepoName(repoName string) (string, error) {
	if repoName == "" {
		return "", errors.New("must specify the repository path")
	}

	if strings.HasPrefix(repoName, "https://") {
		if strings.HasSuffix(repoName, ".git") {
			return repoName, nil
		}
		return repoName + ".git", nil
	}
	if strings.HasPrefix(repoName, "git@") {
		if strings.HasSuffix(repoName, ".git") {
			return repoName, nil
		}
		return repoName + ".git", nil
	}

	if strings.ContainsRune(repoName, '/') {
		return fmt.Sprintf("https://github.com/%s.git", repoName), nil
	}
	return fmt.Sprintf("https://github.com/%s/dotfiles.git", repoName), nil
}
