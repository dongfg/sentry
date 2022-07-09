package internal

import (
	"fmt"
	"log"
	"net/http"
)

func Start(port int64, client *GitClient) {
	// TODO clone
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/webhook" {
			// TODO pull
			return
		}
	})
	addr := fmt.Sprintf(":%d", port)
	log.Println("Listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
