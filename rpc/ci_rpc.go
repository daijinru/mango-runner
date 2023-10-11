package rpc

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/daijinru/mango/mango-cli/runner"
	"github.com/daijinru/mango/mango-cli/utils"
)

var (
  CI_LOCK_NAME = ".running.lock"
)

type CiService struct {
}

type CreatePipelineReply struct {
  Tag string
  Running bool
  Content string
}
type Reply struct {
  Status int8
  Message string
  Data CreatePipelineReply 
}
type CreatePipArgs struct {
  Tag string
  Path string
}

func formatPipMsg(ci *runner.CiClient, msg string) string {
  return fmt.Sprintf("[%s] [%s] 🥭 %s", utils.TimeNow(), ci.Pipeline.Filename, msg)
}
// Path parameter passing that service will switch to the path,
// as the working directory,
// and then performing the tasks by meta-inf/.mango-ci.yaml
func (CiS *CiService) CreatePip(args *CreatePipArgs, reply *Reply) error {
  reply.Status = int8(FailedCreate)

  ciOption := &runner.CiOption{
    Path: args.Path,
    LockName: CI_LOCK_NAME,
    Tag: args.Tag,
  }

  ci := &runner.CiClient{}
  _, err := ci.NewCI(ciOption)
  if err != nil {
    message := fmt.Sprintf("❌ error occured at NewCI: %s", err)
    reply.Message = message
    return nil
  }

  defer ci.Logger.Writer.Close()
  defer ci.Pipeline.File.Close()

  running, err := ci.AreRunningLocally()
  if err != nil {
    ci.Logger.ReportWarn(err.Error())
    reply.Message = formatPipMsg(ci, err.Error())
    return nil
  }
  if running {
    message := "🔒 ci is running, No further operations allowed until it ends"
    reply.Message = formatPipMsg(ci, message)
    ci.Logger.ReportWarn(message)
    return nil
  }

  ok, err := ci.ReadFromYaml()
  if ok {
    ci.Logger.ReportLog("📝 ci completes reading from local yaml")
  } else {
    message := fmt.Sprintf("error occured at ci.ReadFromYaml: %s", err)
    reply.Message = formatPipMsg(ci, message)
    ci.Logger.ReportWarn(err.Error())
    return nil
  }

  err = ci.CreateRunningLocally()
  if err != nil {
    ci.Logger.ReportWarn(err.Error())
    reply.Message = formatPipMsg(ci, "create lock fail")
    return nil
  } else {
    ci.Logger.ReportLog("🔒 create lock file locally success")
  }

  reply.Status = int8(OK)
  reply.Message = formatPipMsg(ci, "new pipeline was successfully launched!")
  reply.Data.Tag = ci.Pipeline.Tag

  // Executiing of the pipelines is time-consuming,
  // do not wait here just let for reponding
  go func() {
    execution := &runner.Execution{
      PrintLine: true,
      Pipeline: ci.Pipeline,
    }
    OuterLoop:
    for stage := ci.Stages.Front(); stage != nil; stage = stage.Next() {
      scripts := stage.Value
      if value, ok := scripts.(*runner.Job); ok {
        ci.Logger.ReportLog("🎯 now running stage: " + value.Stage)
        // fmt.Println(value)
        for _, script := range value.Scripts {
          _, err := execution.RunCommandSplit(script.(string))
          if err != nil {
            ci.Logger.ReportWarn(fmt.Sprintf("❌ has launched stage: [%s], but execution of ci script failed: %s", value.Stage, err))
            ci.Logger.ReportWarn(fmt.Sprintf("sorry 😅, the task was interrupted cause of error occured in stage: [%s], pipelind id: [%s]", value.Stage, ci.Pipeline.Tag))
            break OuterLoop
          }
        }
      }
    }

    err = ci.CompletedRunningTask()
    if err != nil {
      ci.Logger.ReportWarn(fmt.Sprintf("unable ended running pipeline: %s", err))
    }
    if ok {
      ci.Logger.ReportSuccess("✅ finish running task and now release 🔓 the lock")
    }
  }()
  return nil
}

type QueryPipArgs struct {
  Tag string
  DateTime string
  Path string
  Filename string
}

func (Cis *CiService) GetPipStatus (args *QueryPipArgs, reply *Reply) error {
  reply.Status = int8(FailedQuery)
  workspace := &runner.WorkspaceClient{}
  workspace.NewWorkSpaceClient(args.Path)
  lockFilePath, err := url.JoinPath(workspace.CWD, CI_LOCK_NAME)
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  pipelinePath, err := url.JoinPath(workspace.CWD, "./meta-inf/pipelines/")
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  pip, err := runner.NewPipeline(args.Tag, pipelinePath)
  if err != nil {
    reply.Message = err.Error()
    return nil
  }

  existed, err := pip.ArePipelineRunning(lockFilePath)
  reply.Status = int8(OK)
  if err != nil {
    reply.Message = err.Error()
  }
  reply.Data.Running = existed
  return nil
}

// Jobs tasks execution is output to a file, and its calling returns the contents of the file.
func (Cis *CiService) GetPipStdout (args *QueryPipArgs, reply *Reply) error {
  reply.Status = int8(FailedQuery)
  
  workspace := &runner.WorkspaceClient{}
  workspace.NewWorkSpaceClient(args.Path)
  directory, err := url.JoinPath(workspace.CWD, "./meta-inf/pipelines/")
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  files, err := os.ReadDir(directory)
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  for _, file := range files {
    if file.Name() == args.Filename {
      filePath := filepath.Join(directory, file.Name())
      content, err := os.ReadFile(filePath)
      if err != nil {
        reply.Message = err.Error()
        return nil
      }
      reply.Status = int8(OK)
      reply.Data.Content = string(content)
      return nil
    }
  }
  reply.Status = int8(OK)
  reply.Data.Content = "not found stdout the pipeline"
  return nil
}

type PipsPagination struct {
  Total int
  Filenames []string
  Tag string
}
type PipsReply struct {
  Status int8
  Message string
  Data PipsPagination
}
// Gets all pipeline files by the path passing.
func (Cis *CiService) GetPips (args *QueryPipArgs, reply *PipsReply) error {
  reply.Status = int8(FailedQuery)
  workspace := &runner.WorkspaceClient{}
  workspace.NewWorkSpaceClient(args.Path)
  path, err := url.JoinPath(workspace.CWD, "./meta-inf/pipelines/")
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  pip, err := runner.NewPipeline(args.Tag, path)
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  pipFilenames, err := pip.List()
  if err != nil {
    reply.Message = err.Error()
    return nil
  }
  reply.Status = int8(OK)
  reply.Data.Total = len(pipFilenames)
  reply.Data.Filenames = pipFilenames
  reply.Data.Tag = pip.Tag
  return nil
}
