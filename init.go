package main
import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kargirwar/break-timer/types"
	"github.com/kargirwar/break-timer/alarms"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_FILE = "break-timer.log"
var timerCh chan []types.Rule
var logger *lumberjack.Logger

func init() {
	log.Println("init")
	timerCh = make(chan []types.Rule)
	var rules []types.Rule
	go alarms.Start(timerCh, rules)

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
		log.Println(home + "/Library/BreakTimer/" + LOG_FILE);
		return home + "/Library/BreakTimer/" + LOG_FILE
	}

	return LOG_FILE
}

