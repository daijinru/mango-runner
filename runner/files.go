package runner

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"syscall"

	"github.com/daijinru/mango/mango-cli/utils"
)

// It's the working directory client.
type WorkspaceClient struct {
  CWD string `json:"Worksapce"`
  LockFile *LockFile `json:"LockFile"`
}

// Specify a path and use it to initialize current workspace, add CWD to the instance.
func (client *WorkspaceClient) NewWorkSpaceClient(path string) (*WorkspaceClient, error) {
  workspace, err := client.chWorkspace(path)
  client.CWD = workspace
  return client, err
}

func (client *WorkspaceClient) chWorkspace(path string) (string, error) {
  err := os.Chdir(path)
  if err != nil {
    return "", err
  }

  if (!client.PathExists(path)) {
    return "", fmt.Errorf("path not exist: %s", path)
  }

  dir, err := os.Getwd()
  if err != nil {
    return "", fmt.Errorf("no access: %s", path)
  }

  return dir, nil
}

func (client *WorkspaceClient) ListDirectories(path string) ([]string, error) {
  var directories []string

  files, err := os.ReadDir(path)
  if err != nil {
    return nil, err
  }

  for _, file := range files {
    if file.IsDir() {
      directories = append(directories, file.Name())
    }
  }

  return directories, nil
}

func (client *WorkspaceClient) LsFiles(path string) ([]string, error) {
  var out []string

  files, err := os.ReadDir(path)
  if err != nil {
    return nil, err
  }

  for _, file := range files {
      out = append(out, file.Name())
  }

  return out, err
}

func (client *WorkspaceClient) PathExists(path string) bool {
  _, err := os.Stat(path)
  if (err != nil && os.IsNotExist(err)) {
    return false
  }
  return true
}

type LockFile struct {
  Timestamp string
  Name string
  LockFilePath string
}
// Specify an id to initialize the LockFile instance, which is build-in the WorkspaceClient,
// must be executed before using FileLock operations.
func (client *WorkspaceClient) NewLockFile(name string) *WorkspaceClient {
  suffix := ".lock"
  lockFile := &LockFile{
    Name: name,
    Timestamp: utils.TimeNow(),
    LockFilePath: name + suffix,
  }
  client.LockFile = lockFile
  return client
} 

func (client *WorkspaceClient) IfExistsLock() (bool, error) {
  _, err := os.Stat(client.LockFile.LockFilePath)
  if err == nil {
    return true, nil
  } else if os.IsNotExist(err) {
    return false, nil
  } else {
    return false, fmt.Errorf("unable to check file: %v", err)
  }
}

func (client *WorkspaceClient) CreateLockFile() (bool, error) {
  lockFile := client.LockFile
  if _, err := os.Stat(lockFile.LockFilePath); err == nil {
    return true, nil
  } else if os.IsNotExist(err) {
    file, err := os.Create(lockFile.LockFilePath)
    if err != nil {
      return false, fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close()
    return true, nil
  } else {
    return false, fmt.Errorf("unable to check file: %v", err)
  }
}

func (client *WorkspaceClient) DeleteLockFile() error {
  lockFile := client.LockFile
  if _, err := os.Stat(lockFile.LockFilePath); err == nil {
    err := os.Remove(lockFile.LockFilePath)
    if err != nil {
      return fmt.Errorf("lock file cannot be deleted: %v", err)
    }
  }
  return nil
}

var (
  PID_FILE_NAME = ".pid.lock"
)

type Pid struct {
  Pid int
  PidFilePath string
}

type PidOption struct {
  Path string
  // Restart bool
}

func (pid *Pid) ThinClient(option *PidOption) *Pid {
  value, _ := url.JoinPath(option.Path, PID_FILE_NAME)
  pid.PidFilePath = value
  return pid
}

// Check if pid.lock exists,
// if it exists then read PID from old file, check whether the process(by PID) exists,
// (if process not exists) remove the file, and write a new pid.lock.
// if not exists then get a new PID and write it into the new file.
func (pid *Pid) NewPid(option *PidOption) (*Pid, error) {
  value, _ := url.JoinPath(option.Path, PID_FILE_NAME)
  pid.PidFilePath = value
  _, err := os.Stat(pid.PidFilePath)
  if err == nil {
    id, err := pid.ReadPIDFromFile()
    if err != nil {
      return nil, err
    }
    pid.Pid = id
    if pid.ProcessExists() {
      return pid, nil
    } else {
      err := os.Remove(pid.PidFilePath)
      if err != nil {
        return nil, err
      }
      id, err := pid.WritePIDToFile()
      if err != nil {
        return nil, err
      }
      pid.Pid = id
      return pid, nil
    }
  } else {
    id, err := pid.WritePIDToFile()
    if (err != nil) {
      return nil, err
    }
    pid.Pid = id
    return pid, nil
  }
}

func (pid *Pid) WritePIDToFile() (int, error) {
  id := os.Getpid()
  err := os.WriteFile(pid.PidFilePath, []byte(strconv.Itoa(id)), 0644)
  if err != nil {
    return 0, fmt.Errorf("unable write pid file: %v", err)
  }
  return id, nil
}

func (pid *Pid) ReadPIDFromFile() (int, error) {
  bytes, err := os.ReadFile(pid.PidFilePath)
  if err != nil {
    return 0, fmt.Errorf("unable read pid from file: %v", err)
  }
  id, err := strconv.Atoi(string(bytes))
  if err != nil {
    return 0, fmt.Errorf("invalid pid: %v", err)
  }
  return id, nil
}

func (pid *Pid) ProcessExists() bool {
  process, err := os.FindProcess(pid.Pid)
  if err != nil {
    fmt.Printf("unable to get process infor: %v\n", err)
    return false
  }
  err = process.Signal(syscall.Signal(0))
  if err != nil {
    fmt.Printf("the process is not existed: %v\n", err)
    return false
  }
  return true
}

func (pid *Pid) ProcessKill() error {
  id := pid.Pid
  if id == 0 {
    value, err := pid.ReadPIDFromFile()
    if err != nil {
      return fmt.Errorf("unable read PID from file locally: %v", err)
    }
    id = value
  }
  process, err := os.FindProcess(id)
  if err != nil {
    return fmt.Errorf("unable to get process information: %v", err)
  }
  err = process.Signal(syscall.SIGTERM)
  if err != nil {
    return fmt.Errorf("unable to send SIGTERM signal: %v", err)
  }
  _, err = process.Wait()
  if err != nil {
    return fmt.Errorf("unable to wait for process exit: %v", err)
  }
  return nil
}