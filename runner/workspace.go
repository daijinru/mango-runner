package runner

import (
	"os"
	"path/filepath"
)

// Workspace It's the working directory client.
type Workspace struct {
	CWD         string `json:"Workspace"`
	ProjectRoot string `json:"ProjectRoot"`
}

// NewWorkSpace
// The CWD should be <home_dir>/.mango/,
// path should be the project's name, and ProjectRoot will be <home_dir>/mangoes/<project_name>.
// Directory will be created if it doesn't exist.
func NewWorkSpace(path string) (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	absPath := filepath.Join(homeDir, "mangoes", path)
	err = PathExists(absPath)
	if err != nil {
		err := MakePathExists(absPath)
		if err != nil {
			return nil, err
		}
	}

	cwd, err := chWorkspace(filepath.Join(homeDir, ".mango"))
	if err != nil {
		return nil, err
	}
	return &Workspace{
		CWD:         cwd,
		ProjectRoot: absPath,
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
