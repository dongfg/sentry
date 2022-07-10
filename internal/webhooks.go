package internal

import (
	"fmt"
	"github.com/google/go-github/v45/github"
	"io/ioutil"
	"net/http"
)

func gitHubWebhook(w http.ResponseWriter, r *http.Request, client *GitClient) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error reading request body: err=%s\n", err)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		_, _ = fmt.Fprintf(w, "could not parse webhook: err=%s\n", err)
		return
	}

	switch event.(type) {
	case *github.PushEvent:
		_, _ = fmt.Fprintf(w, "git pull\n")
		err = client.Pull(w)
		_, _ = fmt.Fprintf(w, "%v\n", err)
	default:
		return
	}
}
