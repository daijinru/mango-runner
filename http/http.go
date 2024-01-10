package httpService

import (
	"encoding/json"
	"net/http"
	"github.com/daijinru/mango-runner/runner"
)

type CiService struct {
}

type PipelineReply struct {
  Status string `json:"status"`
  Message string `json:"message"`
}

// Path parameter passing that service will switch to the path,
// as the working directory,
// and then performing the tasks by meta-inf/.mango-ci.yaml
func (CiS *CiService) CreatePipeline(w http.ResponseWriter, r *http.Request) {

  runnerArgs := &runner.RunnerArgs{
    Path: r.FormValue("path"),
    Tag: r.FormValue("tag"),
    Callback: r.FormValue("callbackUrl"),
  }
  runner, err := runner.NewRunner(runnerArgs)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  err = runner.Lock.Lock()
  if err != nil {
    runner.Logger.Warn(err.Error())
    message := "🔒 ci is running, No further operations allowed until it ends"
    runner.Logger.Warn(message)
    http.Error(w, message, http.StatusLocked)
    return
  }

  go func() {
    err = runner.Create()
    if err != nil {
      runner.Logger.Warn(err.Error()) 
    }
    runner.Complete()
  }()

  reply := PipelineReply{
    Status: "success",
    Message: runner.Pipeline.Filename,
  }
  jsonData, err := json.Marshal(reply)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}

// Jobs tasks execution is output to a file, and its calling returns the contents of the file.
func (Cis *CiService) ReadPipeline(w http.ResponseWriter, r *http.Request) {
  workspace, err := runner.NewWorkSpace(r.FormValue("path"))
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  pipeline, err := runner.NewPipeline(&runner.PipelineArgs{
    // Tag: r.FormValue("tag"),
    Path: workspace.CWD,
  })
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  reply := PipelineReply{
    Status: "success",
    Message: pipeline.ReadFile(r.FormValue("filename")),
  }
  jsonData, err := json.Marshal(reply)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}

func (Cis *CiService) ReadServiceStatus(w http.ResponseWriter, r *http.Request) {
  reply := PipelineReply{
    Status: "success",
    Message: "true",
  }
  jsonData, err := json.Marshal(reply)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}

type PipelinesReply struct {
  Total int `json:"total"`
  Filenames []string `json:"filenames"`
  Tag string `json:"tag"`
}
// Gets all pipeline files by the path passing.
func (Cis *CiService) ReadPipelines(w http.ResponseWriter, r *http.Request) {
  workspace, err := runner.NewWorkSpace(r.FormValue("path"))
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  pipeline, err := runner.NewPipeline(&runner.PipelineArgs{
    Tag: r.FormValue("tag"),
    Path: workspace.CWD,
  })
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  filenames, err := pipeline.ReadDir()
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  reply := PipelinesReply{
    Total: len(filenames),
    Tag: r.FormValue("tag"),
    Filenames: filenames,
  }
  jsonData, err := json.Marshal(reply)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonData)
}
