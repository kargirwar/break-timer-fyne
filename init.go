package main

import (
	"os"
	"time"
	"path/filepath"
	"runtime"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_FILE = "break-timer.log"
var logger *lumberjack.Logger

func init() {
	//set up logging
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.StampMilli,
	})

	logger = &lumberjack.Logger{
		Filename:   getLogFileName(),
		MaxSize:    1, //megabytes
		MaxBackups: 2,
		MaxAge:     28, //days
		Compress:   true,
	}
	log.SetOutput(logger)

	log.SetLevel(log.DebugLevel)

	timerCh = make(chan []Rule)
	playerCh = make(chan string)
	ch = make(chan *uiCmd)
	ctrls = make(map[string]*TimerCtrl)
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
