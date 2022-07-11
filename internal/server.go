package internal

import (
	"fmt"
	"log"
	"net/http"
)

func Start(port int64, client *GitClient) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
