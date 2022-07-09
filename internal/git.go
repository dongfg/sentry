package internal

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io"
	"os"
)

type GitClient struct {
	rp *git.Repository

	co *git.CloneOptions
}

type GitOptions struct {
	Repo       string
	Username   string
	Password   string
	PrivateKey string
}

func New(o *GitOptions) (client *GitClient) {
	var auth transport.AuthMethod
	var err error
	if o.PrivateKey != "" {
		auth, err = ssh.NewPublicKeysFromFile(o.Username, o.PrivateKey, o.Password)
		if err != nil {
			panic(err)
		}
	} else if o.Password != "" {
		auth = &http.BasicAuth{
			Username: o.Username,
			Password: o.Password,
		}
	}

	return &GitClient{
		co: &git.CloneOptions{
			URL:      o.Repo,
			Auth:     auth,
			Progress: os.Stdout,
		},
	}
}

func (c *GitClient) Clone(path string) {
	r, err := git.PlainClone(path, false, c.co)
	if err != nil {
		panic(err)
	}
	c.rp = r
}

func (c *GitClient) Pull(wr io.Writer) error {
	w, err := c.rp.Worktree()
	if err != nil {
		return err
	}
	return w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   wr,
	})
}
