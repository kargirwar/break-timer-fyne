package main

import (
	_ "fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var PORT = 22890
var VERSION = "VERSION"
const LOG_FILE = "break-timer.log"

var logger *lumberjack.Logger

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.StampMilli,
	})

	logger = &lumberjack.Logger{
		Filename:   getLogFileName(),
		MaxSize:    1, //megabytes
		MaxBackups: 1,
		MaxAge:     1, //days
		Compress:   true,
	}
	log.SetOutput(logger)

	log.SetLevel(log.DebugLevel)
}

func getLogFileName() string {
	//Linux
	if runtime.GOOS == "linux" {
		return "/var/log/break-timer/" + LOG_FILE
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return LOG_FILE
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
	//r.HandleFunc("/about", about).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/set-timers", setTimers).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/stop", stop).Methods(http.MethodGet, http.MethodOptions)

	http.Handle("/", r)

	log.Info("bt-agent Listening at:" + strconv.Itoa(PORT))

	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil); err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}
