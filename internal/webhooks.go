package internal

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

	if client.webhookSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		signature = strings.ReplaceAll(signature, "sha256=", "")

		key := []byte(client.webhookSecret)
		mac := hmac.New(sha256.New, key)
		mac.Write(payload)
		calculated := hex.EncodeToString(mac.Sum(nil))
		if calculated != signature {
			_, _ = fmt.Fprintf(w, "verify signature failed")
			return
		}
	}

	event := r.Header.Get("X-GitHub-Event")
	if event == "push" {
		_, _ = fmt.Fprintf(w, "git pull\n")
		err = client.Pull(w)
		_, _ = fmt.Fprintf(w, "%v\n", err)
	} else {
		_, _ = fmt.Fprintf(w, "[event %s] skip pull\n", event)
	}
}

func codingWebhook(w http.ResponseWriter, r *http.Request, client *GitClient) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error reading request body: err=%s\n", err)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	if client.webhookSecret != "" {
		signature := r.Header.Get("X-Coding-Signature")
		signature = strings.ReplaceAll(signature, "sha1=", "")

		key := []byte(client.webhookSecret)
		mac := hmac.New(sha1.New, key)
		mac.Write(payload)
		calculated := hex.EncodeToString(mac.Sum(nil))
		if calculated != signature {
			_, _ = fmt.Fprintf(w, "verify signature failed")
			return
		}
	}

	event := r.Header.Get("X-Coding-Event")
	if event == "push" {
		_, _ = fmt.Fprintf(w, "git pull\n")
		err = client.Pull(w)
		_, _ = fmt.Fprintf(w, "%v\n", err)
	} else {
		_, _ = fmt.Fprintf(w, "[event %s] skip pull\n", event)
	}
}
