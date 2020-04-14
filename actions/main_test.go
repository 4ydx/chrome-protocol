package actions

import (
	"html/template"
	"net/http"
)

var FirstTimeRun = true

var BrowserPath = "/usr/bin/google-chrome"

func HelloServer(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testing.html")
	t.Execute(w, nil)
}

func LocalServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	go func() {
		if FirstTimeRun {
			http.HandleFunc("/", HelloServer)
		}
		FirstTimeRun = false

		err := srv.ListenAndServe()
		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()
	return srv
}
