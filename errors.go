package got

import "fmt"

type InvalidParamError struct {
	Message string
}

func (e *InvalidParamError) Error() string {
	return fmt.Sprintf("invalid parameter: %s", e.Message)
}

type AlreadyInstalledError struct {
	Path PackagePath
}

func (e *AlreadyInstalledError) Error() string {
	return fmt.Sprintf("already installed: %s", e.Path)
}
