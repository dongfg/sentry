package internal

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	ssh2 "golang.org/x/crypto/ssh"
	"io"
	"os"
)

type GitClient struct {
	rp *git.Repository

	co *git.CloneOptions

	webhookSecret string
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
		auth.(*ssh.PublicKeys).HostKeyCallback = ssh2.InsecureIgnoreHostKey()
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
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		r, err := git.PlainOpen(path)
		if err != nil {
			panic(err)
		}
		c.rp = r
		err = c.Pull(os.Stdout)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			panic(err)
		}
		return
	}
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
		Auth:       c.co.Auth,
	})
}

func (c *GitClient) SetWebhookSecret(webhookSecret string) {
	c.webhookSecret = webhookSecret
}
