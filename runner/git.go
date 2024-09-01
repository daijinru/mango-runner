package runner

import (
	"errors"
	"fmt"
	"github.com/daijinru/mango-runner/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"path/filepath"
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

func (gClient *GitClient) DispatchIfExisted() error {
	logger := gClient.Logger
	repoName := gClient.Args.Name
	repoPath := filepath.Join(gClient.Workspace.ProjectRoot)
	fmt.Println(gClient.Workspace.ProjectRoot)
	r, err := git.PlainOpen(repoPath)
	// if repo not exists
	if errors.Is(err, git.ErrRepositoryNotExists) {
		err = gClient.clone()
		if err != nil {
			return err
		}
		return nil
	}
	// if repo existed, keep continually fetching
	wt, err := r.Worktree()
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to get the worktree of [%s]: [%s]", repoName, err))
		return err
	}
	ref, err := r.Reference(plumbing.NewBranchReferenceName(gClient.Args.Branch), true)
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to find the referece for branch [%s]: [%v]", repoName, err))
		return err
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to checkout to branch [%s]: [%v]", repoName, err))
		return err
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: gClient.Args.User,
			Password: gClient.Args.Pwd,
		},
		Progress: &ProgressWriter{},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to pull the latest changes from [%s] [%s]: [%v]", repoName, gClient.Args.Branch, err))
		return err
	}
	return nil
}

func (gClient *GitClient) clone() error {
	auth := &http.BasicAuth{
		Username: gClient.Args.User,
		Password: gClient.Args.Pwd,
	}
	respond, err := git.PlainClone(gClient.Workspace.ProjectRoot, false, &git.CloneOptions{
		URL:           gClient.Args.Repo,
		Progress:      os.Stdout,
		Auth:          auth,
		ReferenceName: plumbing.NewBranchReferenceName(gClient.Args.Branch),
		SingleBranch:  true,
	})
	if err != nil {
		gClient.Logger.Warn(fmt.Sprintf("❌ clone [%s] error!", gClient.Args.Repo))
		return err
	}
	wt, err := respond.Worktree()
	if err != nil {
		gClient.Logger.Warn(fmt.Sprint("❌ get Worktree err!"))
		return err
	}
	status, err := wt.Status()
	if err != nil {
		return err
	}
	callback := gClient.Args.Callback
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
	gClient.Logger.Info(fmt.Sprintf("cloned message: [%s]", respStr))
	return nil
}
