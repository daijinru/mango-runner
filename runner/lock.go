package runner

import (
	"fmt"
	"os"
	"path/filepath"
  "github.com/gofrs/flock"
)

var (
  PID_FILE_NAME = ".pid.lock"
)

type Lock struct {
  FilePath string
  FileLock *flock.Flock
}

func NewLock(path string) *Lock {
  return &Lock{
    FilePath: filepath.Join(path, "running.lock"),
  }
}

func (lock *Lock) Lock() error {
  _, err := os.OpenFile(lock.FilePath, os.O_CREATE, 0666)
  if err != nil {
    return err
  }
  fileLock := flock.New(lock.FilePath)
  locked, err := fileLock.TryLock()
  if err != nil {
    return fmt.Errorf("failed to lock file: %v", lock.FilePath)
  }
  if locked {
    lock.FileLock = fileLock
  } else {
    return fmt.Errorf("file is already locked by another process")
  }
  return nil
}

func (lock *Lock) Unlock() error {
  defer lock.FileLock.Unlock()
  return nil
}
