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

var VERSION = "VERSION"

const PORT = 22890
const LOG_FILE = "break-timer.log"
const SETTINGS_FILE = "settings.json"
const PLAY = "play"
const STOP = "stop"

var logger *lumberjack.Logger
var timerCh chan string
var playerCh chan string

func init() {
	timerCh = make(chan string)
	playerCh = make(chan string)

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.StampMilli,
	})

	logger = &lumberjack.Logger{
		Filename:   getOsFilePath(LOG_FILE),
		MaxSize:    1, //megabytes
		MaxBackups: 1,
		MaxAge:     1, //days
		Compress:   true,
	}
	log.SetOutput(logger)

	log.SetLevel(log.DebugLevel)
}

func getOsFilePath(f string) string {
	//Linux
	if runtime.GOOS == "linux" {
		return "/var/log/break-timer/" + f
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return f
	}

	//OSX
	if runtime.GOOS == "darwin" {
		return home + "/Library/BreakTimer/" + f
	}

	return f
}

func main() {

	go start()
	go play()

	r := mux.NewRouter()
	r.Use(mw)

	//routes
	//r.HandleFunc("/about", about).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/set-timers", setTimers).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/stop", stop).Methods(http.MethodGet, http.MethodOptions)

	http.Handle("/", r)

	log.Info("bt-agent Listening at:" + strconv.Itoa(PORT))

	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil); err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}
