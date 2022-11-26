package main

import (
	_ "fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kargirwar/break-timer/player"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_FILE = "break-timer.log"
var ALLOW = "http://localhost:8080"
var PORT = 21890

var logger *lumberjack.Logger

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.StampMilli,
	})

	logger = &lumberjack.Logger{
		Filename:   getLogFileName(),
		MaxSize:    10, //megabytes
		MaxBackups: 2,
		MaxAge:     28, //days
		Compress:   true,
	}
	log.SetOutput(logger)

	log.SetLevel(log.DebugLevel)
}

func getLogFileName() string {

	home, err := os.UserHomeDir()
	if err != nil {
		return LOG_FILE
	}

	//Linux
	if runtime.GOOS == "linux" {
		path := home + "/.break-timer"
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return LOG_FILE
		}
		return filepath.Join(path, LOG_FILE)
	}

	//OSX
	if runtime.GOOS == "darwin" {
		return home + "/Library/BreakTimer/" + LOG_FILE
	}

	return LOG_FILE
}

func main() {
	r := mux.NewRouter()

	//middleware
	r.Use(mw)

	//routes
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		player.Play(1, 500)
	}).Methods(http.MethodGet, http.MethodOptions)

	http.Handle("/", r)

	log.Info("break-timer Listening at:" + strconv.Itoa(PORT))

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
