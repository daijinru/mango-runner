package runner

import (
	"fmt"
	"time"
)

// Runner The collections of CI methods, plz call NewCI() for initialization.
type Runner struct {
	Pipeline  *Pipeline
	Lock      *Lock
	Workspace *Workspace
	Logger    *Logger
	Execution *Execution
	Callback  string
}

type RunnerArgs struct {
	// Local directory path, absolute or relative
	Name       string
	CommandStr string
	Callback   string
}

// NewRunner Initialize a CI instance.
func NewRunner(args *RunnerArgs) (*Runner, error) {
	// CWD is forced to be the user root.
	workspace, err := NewWorkSpace(args.Name)
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(workspace.CWD)
	if err != nil {
		return nil, err
	}

	pipeline, err := NewPipeline(&PipelineArgs{
		Tag:  args.Name,
		Path: workspace.CWD,
	})
	if err != nil {
		return nil, err
	}

	execution := &Execution{
		CommandStr: args.CommandStr,
		PrintLine:  true,
		Pipeline:   pipeline,
	}

	return &Runner{
		Workspace: workspace,
		Logger:    logger,
		Pipeline:  pipeline,
		Lock:      NewLock(workspace.CWD),
		Execution: execution,
		Callback:  args.Callback,
	}, nil
}

type RunnerContent struct {
	Status bool
}

func (runner *Runner) Create() error {
	startTime := time.Now().Format("2006-01-02 15:04:05")

	_, err := runner.Execution.RunCommandSplit()
	if err != nil {
		runner.Logger.Warn(fmt.Sprintf(
			"❌ for command: [%s], but execution of ci script failed: %s",
			runner.Execution.CommandStr, err))
		runner.Logger.Warn(fmt.Sprintf(
			"sorry 😅, the task was interrupted cause of error occured in command: [%s], pipelind tag: [%s]",
			runner.Execution.CommandStr, runner.Pipeline.Tag))
		return err
	}

	err = runner.Pipeline.WriteInfo("🥭 running completed!" + "\n")
	if err != nil {
		runner.Logger.Warn(err.Error())
		return err
	}
	// Temporary fix, if it is not unlock, the next task will not be executed after send Callback below.
	runner.Complete()

	if runner.Callback != "" {
		err = runner.Pipeline.Callback(
			runner.Callback,
			"endTime", time.Now().Format("2006-01-02 15:04:05"),
			"startTime", startTime,
			"filename", runner.Pipeline.Filename,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (runner *Runner) Complete() error {
	err := runner.Lock.Unlock()
	if err != nil {
		return err
	}
	defer runner.Pipeline.CloseFile()
	return nil
}
