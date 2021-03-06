package internal

import (
	"fmt"
	"log"
	"net/http"
)

func Start(port int64, client *GitClient) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accept request " + r.URL.Path)
		if r.URL.Path == "/" || r.URL.Path == "" {
			_, _ = fmt.Fprintf(w, "pong")
			return
		}
		if r.URL.Path == "/webhook/github" {
			gitHubWebhook(w, r, client)
			return
		}
		if r.URL.Path == "/webhook/coding" {
			codingWebhook(w, r, client)
			return
		}
	})
	addr := fmt.Sprintf(":%d", port)
	log.Println("Listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
