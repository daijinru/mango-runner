package runner

import (
	"os"
	"path/filepath"
)

// Workspace It's the working directory client.
type Workspace struct {
	CWD  string `json:"Workspace"`
	path string `json:"Path"`
}

// NewWorkSpace
// The currently working directory should be the <home_dir>/mango/path/.mango,
// path should be the project's name, and CWD will be as the mango working directory (/path/.mango).
// Directory will be created if it doesn't exist.
func NewWorkSpace(path string) (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	absPath := filepath.Join(homeDir, "mango", path)
	err = PathExists(absPath)
	if err != nil {
		err := MakePathExists(absPath)
		if err != nil {
			return nil, err
		}
	}

	wd, err := chWorkspace(filepath.Join(absPath, ".mango"))
	if err != nil {
		return nil, err
	}
	return &Workspace{
		CWD:  wd,
		path: absPath,
	}, nil
}

func chWorkspace(path string) (string, error) {
	err := PathExists(path)
	if err != nil {
		err := MakePathExists(path)
		if err != nil {
			return "", err
		}
	}

	err = os.Chdir(path)
	if err != nil {
		return "", err
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

func PathExists(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	return nil
}

func MakePathExists(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
