package httpService

import (
	"encoding/json"
	"github.com/daijinru/mango-runner/runner"
	"net/http"
	"strconv"
)

type CiService struct {
}

type HttpResponse[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

// CreatePipeline Path parameter passing that service will switch to the path,
// as the working directory,
// and then performing the tasks by meta-inf/.mango-ci.yaml
func (CiS *CiService) CreatePipeline(w http.ResponseWriter, r *http.Request) {

	runnerArgs := &runner.RunnerArgs{
		Name:       r.FormValue("name"),
		CommandStr: r.FormValue("command"),
		Callback:   r.FormValue("callbackUrl"),
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

	reply := HttpResponse[string]{
		Status: 200,
		//Message: runner.Pipeline.Filename,
		Data: &runner.Pipeline.Filename,
	}
	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ReadPipeline Jobs tasks execution is output to a file, and its calling returns the contents of the file.
func (Cis *CiService) ReadPipeline(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	workspace, err := runner.NewWorkSpace(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pipeline, err := runner.NewPipeline(&runner.PipelineArgs{
		Tag:  name,
		Path: workspace.CWD,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filename := pipeline.ReadFile(r.FormValue("filename"))
	reply := HttpResponse[string]{
		Status: 200,
		Data:   &filename,
	}
	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (Cis *CiService) ReadServiceStatus(w http.ResponseWriter, r *http.Request) {
	reply := &HttpResponse[any]{
		Status:  200,
		Message: "success",
	}
	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

type PipelinesReply struct {
	Total     int      `json:"total"`
	Filenames []string `json:"filenames"`
	Tag       string   `json:"tag"`
}

// ReadPipelines Gets all pipeline files by the path passing.
func (Cis *CiService) ReadPipelines(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	workspace, err := runner.NewWorkSpace(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pipeline, err := runner.NewPipeline(&runner.PipelineArgs{
		Tag:  name,
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

	pipelinesR := PipelinesReply{
		Total:     len(filenames),
		Tag:       r.FormValue("tag"),
		Filenames: filenames,
	}
	reply := HttpResponse[PipelinesReply]{
		Status: 200,
		Data:   &pipelinesR,
	}

	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (cis *CiService) GitClone(w http.ResponseWriter, r *http.Request) {
	gitClient, err := runner.NewGitClient(&runner.GitClientArgs{
		Name:     r.FormValue("name"),
		Repo:     r.FormValue("repo"),
		Branch:   r.FormValue("branch"),
		User:     r.FormValue("user"),
		Pwd:      r.FormValue("pwd"),
		Callback: r.FormValue("callback"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = gitClient.DispatchIfExisted()
	reply := HttpResponse[any]{}
	if err != nil {
		reply.Status = 400
		reply.Message = err.Error()
	} else {
		reply.Status = 200
		reply.Message = "clone success!"
	}
	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (cis *CiService) ReadMonitor(w http.ResponseWriter, r *http.Request) {
	waitStr := r.FormValue("wait")
	if waitStr == "" {
		waitStr = "5"
	}
	wait, err := strconv.Atoi(waitStr)
	if err != nil {
		wait = 5
	}
	monitor := runner.NewSystemClient(wait)
	states, err := monitor.Read()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	statesByte, err := json.Marshal(states)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	statesStr := string(statesByte)
	reply := HttpResponse[string]{
		Status: 200,
		Data:   &statesStr,
	}
	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
