package runner

import (
	"fmt"
	"github.com/daijinru/mango-runner/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

type GitClient struct {
	Workspace *Workspace
	Logger    *Logger
	Args      *GitClientArgs
}

type GitClientArgs struct {
	Name     string `json:"Name"`
	Repo     string `json:"Repo"`
	Branch   string `json:"Branch"`
	User     string `json:"User"`
	Pwd      string `json:"Pwd"`
	Callback string `json:"callback"`
}

func NewGitClient(args *GitClientArgs) (*GitClient, error) {
	path := args.Name
	workspace, err := NewWorkSpace(path)
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(workspace.CWD)
	if err != nil {
		return nil, err
	}

	return &GitClient{
		Workspace: workspace,
		Logger:    logger,
		Args:      args,
	}, err
}

func (gitClient *GitClient) Clone() error {
	auth := &http.BasicAuth{
		Username: gitClient.Args.User,
		Password: gitClient.Args.Pwd,
	}
	respond, err := git.PlainClone(gitClient.Workspace.ProjectRoot, false, &git.CloneOptions{
		URL:           gitClient.Args.Repo,
		Progress:      os.Stdout,
		Auth:          auth,
		ReferenceName: plumbing.NewBranchReferenceName(gitClient.Args.Branch),
		SingleBranch:  true,
	})
	if err != nil {
		gitClient.Logger.Warn(fmt.Sprintf("❌ clone [%s] error!", gitClient.Args.Repo))
		return err
	}
	wt, err := respond.Worktree()
	if err != nil {
		gitClient.Logger.Warn(fmt.Sprint("❌ get Worktree err!"))
		return err
	}
	status, err := wt.Status()
	if err != nil {
		return err
	}
	callback := gitClient.Args.Callback
	respStr := ""
	if callback != "" {
		newQueries := []string{"git_status", status.String()}
		respStr, err = utils.SendCallbackWithHttp(
			callback,
			newQueries,
		)
		if err != nil {
			return err
		}
	}
	gitClient.Logger.Info(fmt.Sprintf("cloned message: [%s]", respStr))
	return nil
}
