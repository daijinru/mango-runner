package runner

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
  "time"
)

// The collections of CI methods, plz call NewCI() for initialization.
type Runner struct {
  Pipeline *Pipeline
  Lock *Lock
  Workspace *Workspace
  YamlReader *YamlReader
  Logger *Logger
  Callback string
}

type RunnerArgs struct {
  // Local directory path, absolute or relative
  Path string
  Tag string
  Callback string
}

// Initialize a CI instance.
func NewRunner(args *RunnerArgs) (*Runner, error) {
  workspace, err := NewWorkSpace(args.Path)
  if err != nil {
    return nil, err
  }

  yamlReader, err := NewYamlReader(workspace.CWD)
  if err != nil {
    return nil, err
  }
  
  logger, err := NewLogger(workspace.CWD)
  if err != nil {
    return nil, err
  }

  pipeline, err := NewPipeline(&PipelineArgs{
    Tag: args.Tag,
    Path: workspace.CWD,
  })
  if err != nil {
    return nil, err
  }

  return &Runner{
    Workspace: workspace,
    YamlReader: yamlReader,
    Logger: logger,
    Pipeline: pipeline,
    Lock: NewLock(workspace.CWD),
    Callback: args.Callback,
  }, nil
}

type RunnerContent struct {
  Status bool
  
}
func (runner *Runner) ReadStatus() (bool, error) {
  _, err := os.Stat(runner.Lock.FilePath)
  if err != nil && os.IsNotExist(err) {
    return false, nil
  }

  lock, err := runner.Lock.Read()
  if err != nil {
    return false, err
  }
  
  pid, err := strconv.Atoi(lock.PID)
  if err != nil {
    return false, err
  }
  process, err := os.FindProcess(pid)
  if err != nil {
    return false, err
  }
  err = process.Signal(syscall.Signal(0))
  if err == nil {
    return true, nil
  }

  return false, nil
}

func (runner *Runner) Create() error {
  startTime := time.Now().Format("2006-01-02 15:04:05")

  err := runner.YamlReader.Read()
  if err != nil {
    return err
  }

  err = runner.Lock.Write(runner.Pipeline.Tag)
  if err != nil {
    return err
  }

  execution := &Execution{
    PrintLine: true,
    Pipeline: runner.Pipeline,
  }
  
  OuterLoop:
  for stage := runner.YamlReader.Stages.Front(); stage != nil; stage = stage.Next() {
    scripts := stage.Value
    if value, ok := scripts.(*Job); ok {
      runner.Logger.Log("🎯 now running stage: " + value.Stage)
      // fmt.Println(value)
      for _, script := range value.Scripts {
        _, err := execution.RunCommandSplit(script.(string))
        if err != nil {
          runner.Logger.Warn(fmt.Sprintf("❌ has launched stage: [%s], but execution of ci script failed: %s", value.Stage, err))
          runner.Logger.Warn(fmt.Sprintf("sorry 😅, the task was interrupted cause of error occured in stage: [%s], pipelind tag: [%s]", value.Stage, runner.Pipeline.Tag))
          break OuterLoop
        }
      }
    }
  }

  runner.Pipeline.WriteInfo("🥭 running completed!" + "\n")
  err = runner.Pipeline.Callback(
    runner.Callback, 
    "status", "1",
    "endTime", time.Now().Format("2006-01-02 15:04:05"),
    "startTime", startTime,
  )
  if err != nil {
    return err
  }
  err = runner.Complete()

  if err != nil {
    return err
  }
  return nil
}

func (runner *Runner) Complete() error {
  _, err := os.Stat(runner.Lock.FilePath)
  if err != nil && os.IsNotExist(err) {
    return nil
  }
  err = runner.Lock.remove()
  if err != nil {
    return err
  }
  defer runner.Pipeline.CloseFile()
  return nil
}
