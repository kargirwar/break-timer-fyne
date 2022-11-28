package main

import (
	_ "fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var ALLOW = "*"
var PORT = 21890


func main() {
	r := mux.NewRouter()

	//middleware
	r.Use(mw)

	//routes
	r.HandleFunc("/set-timers", func(w http.ResponseWriter, r *http.Request) {
	}).Methods(http.MethodGet, http.MethodOptions)

	http.Handle("/", r)

	log.Info("break-timer Listening at:" + strconv.Itoa(PORT))
	//log.Debug("break-timer Listening at:" + strconv.Itoa(PORT))

	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil); err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}

func mw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", ALLOW)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		w.Header().Set("Access-Control-Allow-Headers", "X-Request-ID")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
