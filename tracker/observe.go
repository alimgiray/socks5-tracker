package tracker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (u *usageTracker) Observe(port string) {
	go func() {
		log.Fatal(http.ListenAndServe(port, nil))
	}()
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		l := u.getLogs()

		bytes, err := json.Marshal(l)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "error: %s", err.Error())
			return
		}

		fmt.Fprintf(w, "%s", string(bytes))
	})
}
