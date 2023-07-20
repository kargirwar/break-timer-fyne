package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kargirwar/golang/utils"
)

const SETTINGS_FILE = "settings-v2.json"
const PLAY = "play"
const STOP = "stop"

type uiCmd struct {
	cmd  string
	data interface{}
}

var ch chan *uiCmd
var gui *fyne.Container
var ctrls map[string]*TimerCtrl
var timerCh chan []Rule
var playerCh chan string

func getSavedRules() []Rule {
	//check if any timers have already been setup
	f := getOsFilePath(SETTINGS_FILE)
	settings, err := ioutil.ReadFile(f)

	if err == nil {
		rules := parseRules(string(settings))
		utils.Dbg(context.Background(), fmt.Sprint(rules))
		return rules
	}

	return []Rule{}
}

func parseRules(jsonstr string) []Rule {
	var rules []Rule
	err := json.Unmarshal([]byte(jsonstr), &rules)
	if err == nil {
		return rules
	}

	return []Rule{}
}

func main() {

	rules := getSavedRules()

	go start(rules)
	go play()
	go uiHandler()

	app := app.New()
	window := app.NewWindow("Break Timer")

	gui = container.NewVBox()
	btns := container.NewHBox(
		widget.NewButton("Add rule", func() {
			utils.Dbg(context.Background(), "Adding new rule")
			ctrl := NewTimerCtrl()
			ctrls[ctrl.uid()] = ctrl
			gui.Add(ctrl.ui(Rule{}))
		}),

		widget.NewButton("Save rules", func() {
			setRules()
		}),
	)

	gui.Add(btns)
	//initialize with rules, if any
	for _, r := range rules {
		ctrl := NewTimerCtrl()
		ctrls[ctrl.uid()] = ctrl
		gui.Add(ctrl.ui(r))
	}

	window.SetContent(gui)
	window.Resize(fyne.NewSize(300, 200))
	window.ShowAndRun()
}

func setRules() {
	var rules []Rule
	for _, v := range ctrls {
		rules = append(rules, v.rule())
	}
	json, _ := json.Marshal(rules)
	ioutil.WriteFile(getOsFilePath(SETTINGS_FILE), json, 0644)

	timerCh <- rules
}

func uiHandler() {
	for {
		select {
		case c := <-ch:
			utils.Dbg(context.Background(), "Cmd received:"+c.cmd)
			id := c.data.(string)
			ctrl, ok := ctrls[id]
			if !ok {
				//log.Println("Unable to remove " + id)
				utils.Dbg(context.Background(), "Unable to remove "+id)
				continue
			}
			gui.Remove(ctrl.ui(Rule{}))
			delete(ctrls, id)
		}
	}
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
		root := home + "/Library/BreakTimer"
		os.MkdirAll(root, os.ModePerm)
		return root + "/" + f
	}

	return f
}
